package core

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/wowsims/cata/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	itemStarshardEdge     = 45620
	itemPillarOfFortitude = 46350
	itemIronmender        = 45271
	itemBookOfBindingWIll = 65133
)

var (
	starshardEdge1 = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemStarshardEdge},
		Slot:  proto.ItemSlot_ItemSlotMainHand,
		Index: 0,
	}
	ironmender = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemIronmender},
		Slot:  proto.ItemSlot_ItemSlotOffHand,
		Index: 1,
	}
	pillarOfFortitude = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemPillarOfFortitude},
		Slot:  proto.ItemSlot_ItemSlotMainHand,
		Index: 2,
	}
	bookOfBindingWill = &itemWithSlot{
		Item:  &proto.ItemSpec{Id: itemBookOfBindingWIll},
		Slot:  proto.ItemSlot_ItemSlotOffHand,
		Index: 3,
	}

	tinyItemDatabase = &proto.SimDatabase{
		Items: []*proto.SimItem{
			{Id: itemStarshardEdge, Type: proto.ItemType_ItemTypeWeapon, HandType: proto.HandType_HandTypeMainHand},
			{Id: itemPillarOfFortitude, Type: proto.ItemType_ItemTypeWeapon, HandType: proto.HandType_HandTypeTwoHand},
			{Id: itemIronmender, Type: proto.ItemType_ItemTypeWeapon, HandType: proto.HandType_HandTypeOffHand},
			{Id: itemBookOfBindingWIll, Type: proto.ItemType_ItemTypeWeapon, HandType: proto.HandType_HandTypeOffHand},
		},
	}
)

func TestIsValidEquipment(t *testing.T) {
	// This is a bit awkward because code everywhere accesses the global database maps. Hopefully
	// this won't mess with any other unit tests that need existing item/gem/enchant databases?
	addToDatabase(tinyItemDatabase)

	for _, tc := range []struct {
		comment       string
		spec          *proto.EquipmentSpec
		isFuryWarrior bool
		want          bool
	}{
		{
			comment:       "simple equipment set with just one mainhand weapon is valid",
			spec:          createEquipmentFromItems(starshardEdge1, bookOfBindingWill),
			isFuryWarrior: false,
			want:          true,
		},
		{
			comment:       "simple equipment set with just one two-handed weapon is valid",
			spec:          createEquipmentFromItems(ironmender),
			isFuryWarrior: false,
			want:          true,
		},
		{
			comment:       "simple equipment set with just one mainhand weapon is not valid",
			spec:          createEquipmentFromItems(starshardEdge1),
			isFuryWarrior: false,
			want:          false,
		},
		{
			comment:       "cannot equip offhand and two-hander if player is not a fury warrior",
			spec:          createEquipmentFromItems(pillarOfFortitude, ironmender),
			isFuryWarrior: false,
			want:          false,
		},
		{
			comment:       "can equip offhand and two-hander if player is a fury warrior",
			spec:          createEquipmentFromItems(ironmender, ironmender),
			isFuryWarrior: true,
			want:          true,
		},
	} {
		if got := isValidEquipment(tc.spec, tc.isFuryWarrior); got != tc.want {
			t.Fatalf("%s: isValidEquipment(%v) = %v, want %v", tc.comment, tc.spec, got, tc.want)
		}
	}
}

func createEquipmentFromItems(items ...*itemWithSlot) *proto.EquipmentSpec {
	spec := &proto.EquipmentSpec{
		Items: make([]*proto.ItemSpec, len(proto.ItemSlot_name)),
	}
	for _, is := range items {
		spec.Items[is.Slot] = is.Item
	}
	for i := range spec.Items {
		if spec.Items[i] == nil {
			spec.Items[i] = &proto.ItemSpec{}
		}
	}
	return spec
}

func TestBulkSim(t *testing.T) {
	t.Skip("TODO: Implement")

	fakeRunSim := func(rsr *proto.RaidSimRequest, progress chan *proto.ProgressMetrics, skipPresim bool, quitChan chan bool) *proto.RaidSimResult {
		return &proto.RaidSimResult{}
	}

	bulk := &bulkSimRunner{
		SingleRaidSimRunner: fakeRunSim,
		Request:             &proto.BulkSimRequest{},
	}

	got, err := bulk.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("BulkSim() returned error: %v", err)
	}

	want := &proto.BulkSimResult{}
	if diff := cmp.Diff(want, got, cmp.Comparer(func(a, b *proto.BulkSimResult) bool {
		return protojson.Format(a) == protojson.Format(b)
	})); diff != "" {
		t.Fatalf("BulkSim() returned diff (-want +got):\n%s", diff)
	}
}

func TestGenerateAllEquipmentSubstitutions(t *testing.T) {
	baseItems := make([]*proto.ItemSpec, len(proto.ItemSlot_name))
	for i := range baseItems {
		baseItems[i] = &proto.ItemSpec{Id: int32(i) + 1000}
	}
	item1 := &proto.ItemSpec{Id: 1}
	item2 := &proto.ItemSpec{Id: 2}
	item3 := &proto.ItemSpec{Id: 1010}
	item4 := &proto.ItemSpec{Id: 4}
	type args struct {
		combinations           bool
		distinctItemSlotCombos []*itemWithSlot
	}
	tests := []struct {
		name string
		args args
		want []*equipmentSubstitution
	}{
		{
			name: "no combos",
			args: args{
				combinations:           true,
				distinctItemSlotCombos: []*itemWithSlot{},
			},
			want: []*equipmentSubstitution{
				{},
			},
		},
		{
			name: "one item",
			args: args{
				combinations: true,
				distinctItemSlotCombos: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotHead},
				},
			},
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotHead},
				}},
			},
		},
		{
			name: "two items",
			args: args{
				combinations: true,
				distinctItemSlotCombos: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotHead},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotShoulder},
				},
			},
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotHead},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotHead},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotShoulder},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotShoulder},
				}},
			},
		},
		{
			name: "ring and trinket",
			args: args{
				combinations: true,
				distinctItemSlotCombos: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				},
			},
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
			},
		},
		{
			name: "two rings and one trinket",
			args: args{
				combinations: true,
				distinctItemSlotCombos: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				},
			},
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
				{Items: []*itemWithSlot{
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
				{Items: []*itemWithSlot{
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
				{Items: []*itemWithSlot{
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
				{Items: []*itemWithSlot{
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item4, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket1},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotTrinket2},
				}},
			},
		},
		{
			name: "special case same itemID",
			args: args{
				combinations: false,
				distinctItemSlotCombos: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item3, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item3, Slot: proto.ItemSlot_ItemSlotFinger2},
				},
			},
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item3, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
			},
		},
		{
			name: "special case finger combo",
			args: args{
				combinations: false,
				distinctItemSlotCombos: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotFinger2},
				},
			},
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger1},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotFinger1},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotFinger2},
				}},
			},
		},

		{
			name: "special case weapon combo",
			args: args{
				combinations: true,
				distinctItemSlotCombos: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotOffHand},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotOffHand},
					{Item: item3, Slot: proto.ItemSlot_ItemSlotMainHand},
				},
			},
			want: []*equipmentSubstitution{
				{},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotMainHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: nil, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotMainHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: nil, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item3, Slot: proto.ItemSlot_ItemSlotMainHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item3, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item1, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item3, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: item2, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item3, Slot: proto.ItemSlot_ItemSlotMainHand},
					{Item: nil, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item1, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: item2, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
				{Items: []*itemWithSlot{
					{Item: nil, Slot: proto.ItemSlot_ItemSlotOffHand},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := generateAllEquipmentSubstitutions(context.Background(), baseItems, tt.args.combinations, tt.args.distinctItemSlotCombos)

			idx := 0
			for got := range results {
				wanted := tt.want[idx]
				if len(got.Items) != len(wanted.Items) {
					t.Errorf("%s generateAllEquipmentSubstitutions(%d) has incorrect number of items, expected: %d, got: %d", tt.name, idx, len(wanted.Items), len(got.Items))
				}
				for itemIdx, item := range got.Items {
					if wanted.Items[itemIdx].Item != nil && wanted.Items[itemIdx].Item.Id != item.Item.Id {
						t.Errorf("%s generateAllEquipmentSubstitutions(%d) has incorrect item in list, expected: %d, got: %d", tt.name, idx, wanted.Items[itemIdx].Item.Id, item.Item.Id)
					}
				}
				idx++
			}
			if idx != len(tt.want) {
				t.Errorf("%s generateAllEquipmentSubstitutions has incorrect number of items, expected: %d, got: %d", tt.name, len(tt.want), idx)
			}
		})
	}
}
