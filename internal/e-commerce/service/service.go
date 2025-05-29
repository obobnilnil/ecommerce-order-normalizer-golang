package service

import (
	"e-commerce/internal/e-commerce/helper"
	"e-commerce/internal/e-commerce/model"
	"e-commerce/internal/e-commerce/repository"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Normalizer interface {
	Normalize(input model.InputOrderItem) ([]model.CleanedOrder, error)
}

type ServicePort interface {
	NormalizeOrderService(input []model.InputOrderItem) ([]model.CleanedOrder, error)
}
type defaultNormalizer struct {
	cleanerMap map[string]string
}
type serviceAdapter struct {
	normalizers map[string]Normalizer // key = platform (Shopee, TikTok, etc.) dependency injection: Platform-specific normalizers, injected via constructor (e.g., Shopee, TikTok)
	defaultNorm Normalizer            // dependency injection: Default normalizer to use when platform is not recognized
	r           repository.RepositoryPort
}

func NewDefaultNormalizer(cleanerMap map[string]string) Normalizer {
	return &defaultNormalizer{cleanerMap: cleanerMap}
}

func NewServiceAdapter(r repository.RepositoryPort, cleanerMap map[string]string) ServicePort {
	defaultNorm := NewDefaultNormalizer(cleanerMap)
	return &serviceAdapter{
		r:           r,
		defaultNorm: defaultNorm,
		normalizers: map[string]Normalizer{
			"shopee": defaultNorm,
			"lazada": defaultNorm,
			// "shopee":  shopeeNormalizer,
			// "tiktok":  tiktokNormalizer,
			// "amazon":  amazonNormalizer,
			// "lazada": lazadaNormalizer,
			// "facebook": facebookNormalizer,
		},
	}
}

func (s *serviceAdapter) NormalizeOrderService(input []model.InputOrderItem) ([]model.CleanedOrder, error) {
	var results []model.CleanedOrder
	orderNo := 1

	for _, in := range input {
		norm := s.defaultNorm
		if n, ok := s.normalizers[in.Channel]; ok {
			norm = n
		}

		out, err := norm.Normalize(in)
		if err != nil {
			// continue
			return nil, fmt.Errorf("normalize failed on item no %d (productId: %s): %w", in.No, in.PlatformProductId, err)
		}

		for i := range out {
			out[i].No = orderNo
			orderNo++
		}

		results = append(results, out...)
	}

	results = mergeDuplicateProducts(results)

	if err := s.r.NormalizeOrderRepository(results); err != nil {
		return nil, err
	}

	return results, nil
}

func (n *defaultNormalizer) Normalize(input model.InputOrderItem) ([]model.CleanedOrder, error) {

	raw := input.PlatformProductId
	baseNo := input.No
	totalPrice := input.TotalPrice

	parts := strings.Split(raw, "/")

	type parsedItem struct {
		part      string
		actualQty int
	}
	var parsedItems []parsedItem
	totalQtyAll := 0

	for _, part := range parts {
		part = cleanPrefix(part)
		qty := input.Qty

		if strings.Contains(part, "*") {
			segments := strings.Split(part, "*")
			part = segments[0]
			if q, err := strconv.Atoi(segments[1]); err == nil {
				qty = q
			}
		}

		parsedItems = append(parsedItems, parsedItem{part, qty})
		totalQtyAll += qty
	}

	var results []model.CleanedOrder
	cleanerMap := map[string]int{}

	for _, item := range parsedItems {
		subparts := strings.Split(item.part, "-")
		if len(subparts) < 3 {
			// return nil, fmt.Errorf("invalid product format: %s", item.part)
			return nil, fmt.Errorf("invalid product format: %s (raw input: %s)", item.part, input.PlatformProductId)
		}

		materialId := strings.Join(subparts[0:2], "-")
		modelId := strings.Join(subparts[2:], "-")
		productId := materialId + "-" + modelId

		unitPrice := totalPrice / float64(totalQtyAll)
		total := unitPrice * float64(item.actualQty)

		results = append(results, model.CleanedOrder{
			No:         baseNo + len(results),
			ProductId:  productId,
			MaterialId: materialId,
			ModelId:    modelId,
			Qty:        item.actualQty,
			UnitPrice:  unitPrice,
			TotalPrice: total,
		})

		cleanerMap[strings.ToUpper(subparts[1])] += item.actualQty
	}

	// WIPING-CLOTH
	results = append(results, model.CleanedOrder{
		No:         baseNo + len(results),
		ProductId:  "WIPING-CLOTH",
		Qty:        totalQtyAll,
		UnitPrice:  0,
		TotalPrice: 0,
	})

	// CLEANNER: Sort keys to make output deterministic
	var textures []string
	for texture := range cleanerMap {
		textures = append(textures, texture)
	}
	sort.Strings(textures)

	for _, texture := range textures {
		qty := cleanerMap[texture]
		cleanerName, ok := n.cleanerMap[texture]
		if !ok {
			cleanerName = texture + "-CLEANNER"
		}
		results = append(results, model.CleanedOrder{
			No:         baseNo + len(results),
			ProductId:  cleanerName,
			Qty:        qty,
			UnitPrice:  0,
			TotalPrice: 0,
		})
	}

	return results, nil
}

func cleanPrefix(pid string) string {
	prefixes := []string{"x2-3&", "--", "-x", "x", "%20", "%", "&", " "}
	for {
		changed := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(pid, prefix) {
				pid = pid[len(prefix):]
				changed = true
				break
			}
		}
		if !changed {
			break
		}
	}
	return pid
}

func mergeDuplicateProducts(orders []model.CleanedOrder) []model.CleanedOrder {
	type productInfo struct {
		item        model.CleanedOrder
		firstSeen   int
		isAccessory bool
	}

	productMap := make(map[string]*productInfo)
	var productOrder []string
	var accessoryOrder []string

	for idx, o := range orders {
		key := o.ProductId
		isAccessory := helper.IsAccessory(o.ProductId)

		if info, exists := productMap[key]; exists {
			info.item.Qty += o.Qty
			info.item.TotalPrice += o.TotalPrice
		} else {
			productMap[key] = &productInfo{
				item:        o,
				firstSeen:   idx,
				isAccessory: isAccessory,
			}

			if isAccessory {
				accessoryOrder = append(accessoryOrder, key)
			} else {
				productOrder = append(productOrder, key)
			}
		}
	}

	sort.SliceStable(productOrder, func(i, j int) bool {
		return productMap[productOrder[i]].firstSeen < productMap[productOrder[j]].firstSeen
	})

	sort.SliceStable(accessoryOrder, func(i, j int) bool {
		return productMap[accessoryOrder[i]].firstSeen < productMap[accessoryOrder[j]].firstSeen
	})

	var result []model.CleanedOrder
	no := 1

	for _, key := range productOrder {
		item := productMap[key].item
		item.No = no
		no++
		result = append(result, item)
	}

	for _, key := range accessoryOrder {
		item := productMap[key].item
		item.No = no
		no++
		result = append(result, item)
	}

	return result
}
