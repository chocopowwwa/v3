package maps

import (
	"testing"
)

type entry struct {
	name string
	unit string
	qty  int
}

func TestAddItem(t *testing.T) {
	tests := []struct {
		name     string
		entry    []entry
		expected bool
	}{
		{
			"Invalid measurement unit",
			[]entry{
				{"pasta", "", 0},
				{"onion", "quarter", 0},
				{"pasta", "pound", 0},
			},
			false,
		},
		{
			"Valid measurement unit",
			[]entry{
				{"peas", "quarter_of_a_dozen", 3},
				{"tomato", "half_of_a_dozen", 6},
				{"chili", "dozen", 12},
				{"cucumber", "small_gross", 120},
				{"potato", "gross", 144},
				{"zucchini", "great_gross", 1728},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, item := range tt.entry {
				ok := AddItem(item.name, item.unit)
				if ok != tt.expected {
					t.Errorf("Expected %t from AddItem, found %t at %v", tt.expected, ok, item.name)
				}

				itemQty, ok := bill[item.name]
				if ok != tt.expected {
					t.Errorf("Could not find item %s in customer bill", item.name)
				}

				if itemQty != item.qty {
					t.Errorf("Expected %s to have quantity %d in customer bill, found %d", item.name, item.qty, itemQty)
				}
			}
		})
	}
}

func TestRemoveItem(t *testing.T) {
	type expectedItem struct {
		name   string
		unit   string
		qty    int
		exists bool
	}

	tests := []struct {
		name     string
		remove   []expectedItem
		expected bool
	}{
		{"Item Not found in bill",
			[]expectedItem{
				{"papaya", "gross", 0, false},
			},
			false,
		},
		{"Invalid measurement unit",
			[]expectedItem{
				{"peas", "pound", 3, true},
				{"tomato", "kilogram", 6, true},
				{"cucumber", "stone", 120, true},
			},
			false,
		},
		{"Resulted qty less than 0",
			[]expectedItem{
				{"peas", "half_of_a_dozen", 3, true},
				{"tomato", "dozen", 6, true},
				{"chili", "small_gross", 12, true},
				{"cucumber", "gross", 120, true},
				{"potato", "great_gross", 144, true},
			},
			false,
		},
		{"Should delete the item if 0",
			[]expectedItem{
				{"peas", "quarter_of_a_dozen", 0, false},
				{"tomato", "half_of_a_dozen", 0, false},
				{"chili", "dozen", 0, false},
				{"cucumber", "small_gross", 0, false},
				{"potato", "gross", 0, false},
				{"zucchini", "great_gross", 0, false},
			},
			true,
		},
		{"Should reduce the qty",
			[]expectedItem{
				{"chili", "half_of_a_dozen", 6, true},
				{"cucumber", "dozen", 108, true},
				{"zucchini", "gross", 1584, true},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupInitialBillData()
			for _, item := range tt.remove {
				ok := RemoveItem(item.name, item.unit)
				if ok != tt.expected {
					t.Errorf("Expected %t from RemoveItem, found %t at %v", tt.expected, ok, item.name)
				}

				itemQty, ok := bill[item.name]
				if ok != item.exists {
					t.Errorf("Could not find item %s in customer bill", item.name)
				}
				if itemQty != item.qty {
					t.Errorf("Expected %s to have quantity %d in customer bill, found %d", item.name, item.qty, itemQty)
				}
			}
		})
	}
}

func TestCheckout(t *testing.T) {
	// Success, zero out the  bill
	t.Run("Should reset customerbill after checkout", func(t *testing.T) {
		setupInitialBillData()
		Checkout()

		if len(bill) != 0 {
			t.Error("Customer bill must be empty after checkout")
		}
	})
}

func TestGetItem(t *testing.T) {
	type expectedItem struct {
		name     string
		expected bool
		qty      int
	}

	test := []struct {
		name    string
		getItem []expectedItem
	}{
		{
			"Item Not found in bill",
			[]expectedItem{
				{"zuccini", false, 0},
			},
		},
		{
			"Success",
			[]expectedItem{
				{"peas", true, 3},
				{"tomato", true, 6},
				{"chili", true, 12},
				{"cucumber", true, 120},
				{"potato", true, 144},
				{"zucchini", true, 1728},
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			for _, item := range tt.getItem {
				setupInitialBillData()
				itemQty, ok := GetItem(item.name)

				if ok != item.expected {
					t.Errorf("Could not find item %s in customer bill", item.name)
				}

				if itemQty != item.qty {
					t.Errorf("Expected %s to have quantity %d in customer bill, found %d", item.name, item.qty, itemQty)
				}
			}
		})
	}
}

func setupInitialBillData() {
	bill = map[string]int{}
	bill["peas"] = 3
	bill["tomato"] = 6
	bill["chili"] = 12
	bill["cucumber"] = 120
	bill["potato"] = 144
	bill["zucchini"] = 1728
}