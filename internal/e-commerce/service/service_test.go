// internal/e-commerce/service/service_test.go
package service

import (
	"e-commerce/internal/e-commerce/model"
	"e-commerce/internal/e-commerce/repository"
	"reflect"
	"testing"
)

// func TestNormalize_Case1(t *testing.T) {
// 	cleanerMap := map[string]string{
// 		"CLEAR": "CLEAR-CLEANNER",
// 	}

// 	n := NewDefaultNormalizer(cleanerMap)

// 	input := model.InputOrder{
// 		No:                1,
// 		PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
// 		Qty:               2,
// 		UnitPrice:         50,
// 		TotalPrice:        100,
// 	}

// 	expected := []model.CleanedOrder{
// 		{
// 			No:         1,
// 			ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
// 			MaterialId: "FG0A-CLEAR",
// 			ModelId:    "IPHONE16PROMAX",
// 			Qty:        2,
// 			UnitPrice:  50.0,
// 			TotalPrice: 100.0,
// 		},
// 		{
// 			No:         2,
// 			ProductId:  "WIPING-CLOTH",
// 			Qty:        2,
// 			UnitPrice:  0.0,
// 			TotalPrice: 0.0,
// 		},
// 		{
// 			No:         3,
// 			ProductId:  "CLEAR-CLEANNER",
// 			Qty:        2,
// 			UnitPrice:  0.0,
// 			TotalPrice: 0.0,
// 		},
// 	}

// 	actual, err := n.Normalize(input)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if !reflect.DeepEqual(actual, expected) {
// 		t.Errorf("expected:\n%+v\nbut got:\n%+v", expected, actual)
// 	}
// }

func TestNormalizeOrder(t *testing.T) {
	cleanerMap := map[string]string{
		"CLEAR":   "CLEAR-CLEANNER",
		"MATTE":   "MATTE-CLEANNER",
		"PRIVACY": "PRIVACY-CLEANNER",
	}

	n := NewDefaultNormalizer(cleanerMap)

	tests := []struct {
		name     string
		input    model.InputOrderItem
		expected []model.CleanedOrder
	}{
		{
			name: "Case1: Only one product",
			input: model.InputOrderItem{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50,
				TotalPrice:        100,
			},
			expected: []model.CleanedOrder{
				{No: 1, ProductId: "FG0A-CLEAR-IPHONE16PROMAX", MaterialId: "FG0A-CLEAR", ModelId: "IPHONE16PROMAX", Qty: 2, UnitPrice: 50, TotalPrice: 100},
				{No: 2, ProductId: "WIPING-CLOTH", Qty: 2, UnitPrice: 0, TotalPrice: 0},
				{No: 3, ProductId: "CLEAR-CLEANNER", Qty: 2, UnitPrice: 0, TotalPrice: 0},
			},
		},
		{
			name: "Case 2: Wrong prefix",
			input: model.InputOrderItem{
				No:                1,
				PlatformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
				Qty:               2,
				UnitPrice:         50,
				TotalPrice:        100,
			},
			expected: []model.CleanedOrder{
				{No: 1, ProductId: "FG0A-CLEAR-IPHONE16PROMAX", MaterialId: "FG0A-CLEAR", ModelId: "IPHONE16PROMAX", Qty: 2, UnitPrice: 50, TotalPrice: 100},
				{No: 2, ProductId: "WIPING-CLOTH", Qty: 2, UnitPrice: 0, TotalPrice: 0},
				{No: 3, ProductId: "CLEAR-CLEANNER", Qty: 2, UnitPrice: 0, TotalPrice: 0},
			},
		},
		{
			name: "Case 3: Wrong prefix and *quantity",
			input: model.InputOrderItem{
				No:                1,
				PlatformProductId: "x2-3&FG0A-MATTE-IPHONE16PROMAX*3",
				Qty:               1,
				UnitPrice:         90,
				TotalPrice:        90,
			},
			expected: []model.CleanedOrder{
				{No: 1, ProductId: "FG0A-MATTE-IPHONE16PROMAX", MaterialId: "FG0A-MATTE", ModelId: "IPHONE16PROMAX", Qty: 3, UnitPrice: 30.0, TotalPrice: 90.0},
				{No: 2, ProductId: "WIPING-CLOTH", Qty: 3, UnitPrice: 0.0, TotalPrice: 0.0},
				{No: 3, ProductId: "MATTE-CLEANNER", Qty: 3, UnitPrice: 0.0, TotalPrice: 0.0},
			},
		},
		{
			name: "Case 4: One bundle product with wrong prefix and /",
			input: model.InputOrderItem{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
				Qty:               1,
				UnitPrice:         80,
				TotalPrice:        80,
			},
			expected: []model.CleanedOrder{
				{No: 1, ProductId: "FG0A-CLEAR-OPPOA3", MaterialId: "FG0A-CLEAR", ModelId: "OPPOA3", Qty: 1, UnitPrice: 40, TotalPrice: 40},
				{No: 2, ProductId: "FG0A-CLEAR-OPPOA3-B", MaterialId: "FG0A-CLEAR", ModelId: "OPPOA3-B", Qty: 1, UnitPrice: 40, TotalPrice: 40},
				{No: 3, ProductId: "WIPING-CLOTH", Qty: 2, UnitPrice: 0, TotalPrice: 0},
				{No: 4, ProductId: "CLEAR-CLEANNER", Qty: 2, UnitPrice: 0, TotalPrice: 0},
			},
		},
		{
			name: "Case 5: Wrong prefix with three bundled products",
			input: model.InputOrderItem{
				No:                1,
				PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3",
				Qty:               1,
				UnitPrice:         120,
				TotalPrice:        120,
			},
			expected: []model.CleanedOrder{
				{No: 1, ProductId: "FG0A-CLEAR-OPPOA3", MaterialId: "FG0A-CLEAR", ModelId: "OPPOA3", Qty: 1, UnitPrice: 40, TotalPrice: 40},
				{No: 2, ProductId: "FG0A-CLEAR-OPPOA3-B", MaterialId: "FG0A-CLEAR", ModelId: "OPPOA3-B", Qty: 1, UnitPrice: 40, TotalPrice: 40},
				{No: 3, ProductId: "FG0A-MATTE-OPPOA3", MaterialId: "FG0A-MATTE", ModelId: "OPPOA3", Qty: 1, UnitPrice: 40, TotalPrice: 40},
				{No: 4, ProductId: "WIPING-CLOTH", Qty: 3, UnitPrice: 0, TotalPrice: 0},
				{No: 5, ProductId: "CLEAR-CLEANNER", Qty: 2, UnitPrice: 0, TotalPrice: 0},
				{No: 6, ProductId: "MATTE-CLEANNER", Qty: 1, UnitPrice: 0, TotalPrice: 0},
			},
		},
		{
			name: "Case 6: Wrong prefix + bundle + *quantity",
			input: model.InputOrderItem{
				No:                1,
				PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
				Qty:               1,
				UnitPrice:         120,
				TotalPrice:        120,
			},
			expected: []model.CleanedOrder{
				{No: 1, ProductId: "FG0A-CLEAR-OPPOA3", MaterialId: "FG0A-CLEAR", ModelId: "OPPOA3", Qty: 2, UnitPrice: 40, TotalPrice: 80},
				{No: 2, ProductId: "FG0A-MATTE-OPPOA3", MaterialId: "FG0A-MATTE", ModelId: "OPPOA3", Qty: 1, UnitPrice: 40, TotalPrice: 40},
				{No: 3, ProductId: "WIPING-CLOTH", Qty: 3, UnitPrice: 0, TotalPrice: 0},
				{No: 4, ProductId: "CLEAR-CLEANNER", Qty: 2, UnitPrice: 0, TotalPrice: 0},
				{No: 5, ProductId: "MATTE-CLEANNER", Qty: 1, UnitPrice: 0, TotalPrice: 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := n.Normalize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected: %+v\n got: %+v", tt.expected, got)
			}
		})
	}
	t.Run("Case 7: multiple items", func(t *testing.T) {
		repo := repository.NewRepositoryAdapter(nil)
		s := NewServiceAdapter(repo, cleanerMap)

		input := []model.InputOrderItem{
			{
				No:                1,
				PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
				Qty:               1,
				UnitPrice:         160,
				TotalPrice:        160,
			},
			{
				No:                2,
				PlatformProductId: "FG0A-PRIVACY-IPHONE16PROMAX",
				Qty:               1,
				UnitPrice:         50,
				TotalPrice:        50,
			},
		}

		expected := []model.CleanedOrder{
			{No: 1, ProductId: "FG0A-CLEAR-OPPOA3", MaterialId: "FG0A-CLEAR", ModelId: "OPPOA3", Qty: 2, UnitPrice: 40, TotalPrice: 80},
			{No: 2, ProductId: "FG0A-MATTE-OPPOA3", MaterialId: "FG0A-MATTE", ModelId: "OPPOA3", Qty: 2, UnitPrice: 40, TotalPrice: 80},
			{No: 3, ProductId: "FG0A-PRIVACY-IPHONE16PROMAX", MaterialId: "FG0A-PRIVACY", ModelId: "IPHONE16PROMAX", Qty: 1, UnitPrice: 50, TotalPrice: 50},
			{No: 4, ProductId: "WIPING-CLOTH", Qty: 5, UnitPrice: 0, TotalPrice: 0},
			{No: 5, ProductId: "CLEAR-CLEANNER", Qty: 2, UnitPrice: 0, TotalPrice: 0},
			{No: 6, ProductId: "MATTE-CLEANNER", Qty: 2, UnitPrice: 0, TotalPrice: 0},
			{No: 7, ProductId: "PRIVACY-CLEANNER", Qty: 1, UnitPrice: 0, TotalPrice: 0},
		}

		got, err := s.NormalizeOrderService(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected:\n%+v\nbut got:\n%+v", expected, got)
		}
	})
}
