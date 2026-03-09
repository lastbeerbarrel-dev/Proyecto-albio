package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ResourceAmount struct {
	ResourceID string `json:"resource_id"`
	Count      int    `json:"count"`
}

type ItemMetadata struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Tier           int              `json:"tier"`
	Category       string           `json:"category"`
	SubCategory    string           `json:"sub_category,omitempty"`
	BaseIP         int              `json:"base_ip"`
	ReferencePrice int              `json:"reference_price"`
	HighVolatility bool             `json:"high_volatility,omitempty"`
	Recipe         []ResourceAmount `json:"recipe,omitempty"`
}

func main() {
	items := make(map[string]ItemMetadata)

	resources := map[string]string{
		"METALBAR":   "Lingote de Metal",
		"LEATHER":    "Cuero",
		"CLOTH":      "Tela",
		"PLANKS":     "Tablas",
		"STONEBLOCK": "Bloque de Piedra",
	}

	for t := 4; t <= 8; t++ {
		for id, name := range resources {
			for e := 0; e <= 4; e++ {
				suffix, itemID := "", fmt.Sprintf("T%d_%s", t, id)
				if e > 0 {
					suffix, itemID = fmt.Sprintf(".%d", e), fmt.Sprintf("%s@%d", itemID, e)
				}
				items[itemID] = ItemMetadata{ID: itemID, Name: fmt.Sprintf("%s (T%d%s)", name, t, suffix), Tier: t, Category: "Resource", SubCategory: "Recurso"}
			}
		}
	}

	weaponFamilies := []struct {
		ID       string
		Name     string
		SubCat   string
		C1, C2   int
		R1, R2   string
		Variants []struct {
			ID             string
			Name           string
			HighVolatility bool
		}
	}{
		{"SWORD", "Espada", "Espadas", 16, 8, "METALBAR", "LEATHER", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_SWORD", "Espada Ancha", false},
			{"2H_CLAYMORE", "Claymore", false},
			{"2H_DUALSWORD", "Espadas Dobles", false},
			{"MAIN_SWORD_KEEPER", "Hoja de Clarent", true},
			{"2H_SCIMITAR_MORGANA", "Espada Tallada", true},
			{"2H_DUALSWORD_HELL", "Galatinas", true},
			{"2H_CLAYMORE_AVALON", "Rey", true},
		}},
		{"AXE", "Hacha", "Hachas", 12, 12, "METALBAR", "PLANKS", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_AXE", "Hacha de Batalla", false},
			{"2H_AXE", "Gran Hacha", false},
			{"2H_HALBERD", "Alabarda", false},
			{"2H_AXE_KEEPER", "Carroña", true},
			{"2H_SCYTHE_HELL", "Guadaña Infernal", true},
			{"2H_BEARPAWS", "Patas de Oso", true},
			{"2H_AXE_AVALON", "Rompeteinos", true},
		}},
		{"MACE", "Maza", "Mazas", 16, 8, "METALBAR", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_MACE", "Maza", false},
			{"2H_MACE", "Maza Pesada", false},
			{"2H_FLAIL", "Lucero del Alba", false},
			{"2H_ROCKMACE_KEEPER", "Maza de Roca", true},
			{"2H_MACE_HELL", "Camlann", true},
			{"MAIN_MACE_HELL", "Incubus", true},
			{"2H_MACE_AVALON", "Guardián del Juramento", true},
		}},
		{"HAMMER", "Martillo", "Martillos", 16, 8, "METALBAR", "PLANKS", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_HAMMER", "Martillo", false},
			{"2H_POLEHAMMER", "Martillo de Polo", false},
			{"2H_HAMMER", "Martillo Grande", false},
			{"2H_HAMMER_KEEPER", "Martillo de Tumba", true},
			{"2H_HAMMER_HELL", "Martillo de la Forja", true},
			{"2H_DUALHAMMER_HELL", "Mano de Justicia", true},
			{"2H_HAMMER_AVALON", "Vigilante", true},
		}},
		{"WARGLOVE", "Guantes de Guerra", "Guantes de Guerra", 12, 12, "METALBAR", "LEATHER", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_WARGLOVE", "Guantes de Pelea", false},
			{"2H_WARGLOVE", "Guantes de Batalla", false},
			{"2H_WARGLOVE_SPIKE", "Guantes de Púas", false},
			{"2H_WARGLOVE_KEEPER", "Ursine", true},
			{"2H_WARGLOVE_HELL", "Cuervo", true},
			{"2H_WARGLOVE_MORGANA", "Hellfire", true},
			{"2H_WARGLOVE_AVALON", "Avalona", true},
		}},
		{"BOW", "Arco", "Arcos", 32, 0, "PLANKS", "", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_BOW", "Arco", false},
			{"2H_LONGBOW", "Arco Largo", false},
			{"2H_WARBOW", "Arco de Guerra", false},
			{"2H_BOW_KEEPER", "Susurrante", true},
			{"2H_BOW_HELL", "Arco de Badon", true},
			{"2H_BOW_MORGANA", "Perforanieblas", true},
			{"2H_BOW_AVALON", "Lucent Hawk", true},
		}},
		{"DAGGER", "Daga", "Dagas", 12, 12, "LEATHER", "METALBAR", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_DAGGER", "Daga", false},
			{"2H_DAGGERPAIR", "Dagas en Par", false},
			{"2H_CLAWSSPELL", "Garras", false},
			{"MAIN_DAGGER_KEEPER", "Colmillos", true},
			{"2H_DAGGER_HELL", "Sedienta de Sangre", true},
			{"2H_DAGGER_MORGANA", "Segadora de Almas", true},
			{"2H_DAGGER_AVALON", "Deathgivers", true},
		}},
		{"SPEAR", "Lanza", "Lanzas", 16, 8, "PLANKS", "METALBAR", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_SPEAR", "Lanza", false},
			{"2H_SPEAR", "Pica", false},
			{"2H_GLAIVE", "Guadaña (Glaive)", false},
			{"MAIN_SPEAR_KEEPER", "Garzas", true},
			{"2H_HARPOON", "Espontón", true},
			{"2H_SPEAR_HELL", "Trinidad", true},
			{"2H_SPEAR_AVALON", "Amanecer", true},
		}},
		{"QUARTERSTAFF", "Bastón", "Quarterstaffs", 16, 8, "PLANKS", "METALBAR", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_QUARTERSTAFF", "Bastón de Hierro", false},
			{"2H_QUARTERSTAFF", "Bastón de Monje Negro", false},
			{"2H_IRONCLADSTAFF", "Grial", false},
			{"2H_QUARTERSTAFF_KEEPER", "Equilibrio", true},
			{"2H_QUARTERSTAFF_HELL", "Almas", true},
			{"2H_QUARTERSTAFF_MORGANA", "Buscador de Dobles", true},
			{"2H_QUARTERSTAFF_AVALON", "Cristal", true},
		}},
		{"CROSSBOW", "Ballesta", "Ballestas", 20, 12, "PLANKS", "METALBAR", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_CROSSBOW", "Ballesta Ligera", false},
			{"2H_CROSSBOW", "Ballesta Pesada", false},
			{"2H_CROSSBOW_LARGE", "Ballesta Repetidora", false},
			{"MAIN_CROSSBOW_KEEPER", "Tirador de Saetas", true},
			{"2H_REPEATINGCROSSBOW_HELL", "Ballesta de Asedio", true},
			{"2H_CROSSBOW_MORGANA", "Ballesta de Cristal", true},
			{"2H_CROSSBOW_AVALON", "Avaloniana", true},
		}},
		{"FIRESTAFF", "Bastón de Fuego", "Bastón de Fuego", 16, 8, "PLANKS", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_FIRESTAFF", "Bastón de Fuego", false},
			{"2H_FIRESTAFF", "Gran Bastón de Fuego", false},
			{"2H_INFERNOSTAFF", "Bastón Infernal", false},
			{"2H_FIRESTAFF_KEEPER", "Bastón Salvaje", true},
			{"2H_FIRESTAFF_HELL", "Sulfúrico", true},
			{"2H_FIRESTAFF_MORGANA", "Ardiente", true},
			{"2H_FIRESTAFF_AVALON", "Llamas", true},
		}},
		{"HOLYSTAFF", "Bastón Sagrado", "Bastón Sagrado", 16, 8, "PLANKS", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_HOLYSTAFF", "Bastón Sagrado", false},
			{"2H_HOLYSTAFF", "Gran Bastón Sagrado", false},
			{"2H_DIVINESTAFF", "Bastón Divino", false},
			{"MAIN_HOLYSTAFF_KEEPER", "Vida", true},
			{"2H_HOLYSTAFF_HELL", "Caído", true},
			{"2H_HOLYSTAFF_MORGANA", "Redención", true},
			{"2H_HOLYSTAFF_AVALON", "Hallowfall", true},
		}},
		{"ARCANESTAFF", "Bastón Arcano", "Bastón Arcano", 16, 8, "PLANKS", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_ARCANESTAFF", "Bastón Arcano", false},
			{"2H_ARCANESTAFF", "Grande", false},
			{"2H_ENIGMATICSTAFF", "Enigmático", false},
			{"MAIN_ARCANESTAFF_KEEPER", "Bruja", true},
			{"2H_ARCANESTAFF_HELL", "Ocultismo", true},
			{"2H_ARCANESTAFF_MORGANA", "Malevolencia", true},
			{"2H_ARCANESTAFF_AVALON", "Infinitum", true},
		}},
		{"FROSTSTAFF", "Bastón de Hielo", "Bastón de Hielo", 16, 8, "PLANKS", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_FROSTSTAFF", "Bastón de Hielo", false},
			{"2H_FROSTSTAFF", "Grande", false},
			{"2H_GLACIALSTAFF", "Glacial", false},
			{"MAIN_FROSTSTAFF_KEEPER", "Escarcha", true},
			{"2H_FROSTSTAFF_HELL", "Hoary", true},
			{"2H_FROSTSTAFF_MORGANA", "Permafrost", true},
			{"2H_FROSTSTAFF_AVALON", "Chill", true},
		}},
		{"NATURESTAFF", "Bastón Natural", "Bastón Natural", 16, 8, "PLANKS", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_NATURESTAFF", "Bastón Natural", false},
			{"2H_NATURESTAFF", "Grande", false},
			{"2H_WILDSTAFF", "Salvaje", false},
			{"MAIN_NATURESTAFF_KEEPER", "Druida", true},
			{"2H_NATURESTAFF_HELL", "Plaga", true},
			{"2H_NATURESTAFF_MORGANA", "Ironroot", true},
			{"2H_NATURESTAFF_AVALON", "Espinas", true},
		}},
		{"CURSESTAFF", "Bastón Maldito", "Bastón Maldito", 16, 8, "PLANKS", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_CURSESTAFF", "Bastón Maldito", false},
			{"2H_CURSESTAFF", "Grande", false},
			{"2H_DEMONICSTAFF", "Demoníaco", false},
			{"MAIN_CURSESTAFF_KEEPER", "Calavera", true},
			{"2H_CURSESTAFF_HELL", "Vida", true},
			{"2H_CURSESTAFF_MORGANA", "Sombra", true},
			{"2H_CURSESTAFF_AVALON", "Cristal", true},
		}},
		{"SHAPESHIFTER", "Bastón Cambiaformas", "Bastón Cambiaformas", 16, 8, "PLANKS", "CLOTH", []struct {
			ID             string
			Name           string
			HighVolatility bool
		}{
			{"MAIN_SHAPESHIFTER_PANTHER", "Pantera", false},
			{"2H_SHAPESHIFTER_BEAR", "Oso", false},
			{"2H_SHAPESHIFTER_EAGLE", "Águila", false},
			{"2H_SHAPESHIFTER_TREANT", "Sylvian", true},
			{"2H_SHAPESHIFTER_GOLEM", "Golem", true},
			{"2H_SHAPESHIFTER_WEREWOLF", "Luna de Sangre", true},
			{"2H_SHAPESHIFTER_AVALON", "Cristal", true},
		}},
	}

	for t := 4; t <= 8; t++ {
		for _, f := range weaponFamilies {
			for _, v := range f.Variants {
				for e := 0; e <= 4; e++ {
					suffix, itemID := "", fmt.Sprintf("T%d_%s", t, v.ID)
					if e > 0 {
						suffix, itemID = fmt.Sprintf(".%d", e), fmt.Sprintf("%s@%d", itemID, e)
					}
					recipe := []ResourceAmount{{ResourceID: fmt.Sprintf("T%d_%s", t, f.R1), Count: f.C1}}
					if f.R2 != "" {
						recipe = append(recipe, ResourceAmount{ResourceID: fmt.Sprintf("T%d_%s", t, f.R2), Count: f.C2})
					}
					items[itemID] = ItemMetadata{
						ID:             itemID,
						Name:           fmt.Sprintf("%s (T%d%s)", v.Name, t, suffix),
						Tier:           t,
						Category:       "Weapon",
						SubCategory:    f.SubCat,
						BaseIP:         (300 + t*100) + e*100,
						HighVolatility: v.HighVolatility,
						Recipe:         recipe,
					}
				}
			}
		}
	}

	armorFamilies := []struct{ ID, Name, Res string }{{"PLATE_SET1", "Soldado", "METALBAR"}, {"PLATE_SET2", "Caballero", "METALBAR"}, {"LEATHER_SET1", "Mercenario", "LEATHER"}, {"LEATHER_SET2", "Cazador", "LEATHER"}, {"CLOTH_SET1", "Erudito", "CLOTH"}, {"CLOTH_SET2", "Clerigo", "CLOTH"}}
	for t := 4; t <= 8; t++ {
		for _, f := range armorFamilies {
			parts := []struct {
				ID, Name string
				Count    int
			}{{"HEAD", "Casco", 8}, {"ARMOR", "Peto", 16}, {"SHOES", "Botas", 8}}
			for _, p := range parts {
				for e := 0; e <= 4; e++ {
					suffix, itemID := "", fmt.Sprintf("T%d_%s_%s", t, p.ID, f.ID)
					if e > 0 {
						suffix, itemID = fmt.Sprintf(".%d", e), fmt.Sprintf("%s@%d", itemID, e)
					}
					items[itemID] = ItemMetadata{ID: itemID, Name: fmt.Sprintf("%s de %s (T%d%s)", p.Name, f.Name, t, suffix), Tier: t, Category: "Armor", SubCategory: f.Name, BaseIP: (300 + t*100) + e*100, Recipe: []ResourceAmount{{ResourceID: fmt.Sprintf("T%d_%s", t, f.Res), Count: p.Count}}}
				}
			}
		}
	}

	for t := 4; t <= 8; t++ {
		items[fmt.Sprintf("T%d_BAG", t)] = ItemMetadata{ID: fmt.Sprintf("T%d_BAG", t), Name: fmt.Sprintf("Bolso (T%d)", t), Tier: t, Category: "Accessory", SubCategory: "Bolso", Recipe: []ResourceAmount{{ResourceID: fmt.Sprintf("T%d_CLOTH", t), Count: 8}, {ResourceID: fmt.Sprintf("T%d_LEATHER", t), Count: 8}}}
		items[fmt.Sprintf("T%d_CAPE", t)] = ItemMetadata{ID: fmt.Sprintf("T%d_CAPE", t), Name: fmt.Sprintf("Capa (T%d)", t), Tier: t, Category: "Accessory", SubCategory: "Capa", Recipe: []ResourceAmount{{ResourceID: fmt.Sprintf("T%d_CLOTH", t), Count: 4}}}
		items[fmt.Sprintf("T%d_MOUNT_HORSE", t)] = ItemMetadata{ID: fmt.Sprintf("T%d_MOUNT_HORSE", t), Name: fmt.Sprintf("Caballo (T%d)", t), Tier: t, Category: "Mount", SubCategory: "Caballo"}
		items[fmt.Sprintf("T%d_MOUNT_OX", t)] = ItemMetadata{ID: fmt.Sprintf("T%d_MOUNT_OX", t), Name: fmt.Sprintf("Buey (T%d)", t), Tier: t, Category: "Mount", SubCategory: "Buey"}
	}
	items["T8_MEAL_STEW"] = ItemMetadata{ID: "T8_MEAL_STEW", Name: "Estofado de buey (T8)", Tier: 8, Category: "Consumable", SubCategory: "Comida"}
	items["T7_POTION_HEAL"] = ItemMetadata{ID: "T7_POTION_HEAL", Name: "Poción de sanación (T7)", Tier: 7, Category: "Consumable", SubCategory: "Poción"}

	// Rare/Luxury Items (Cosmetics/Mounts)
	luxuryItems := []struct{ ID, Name, SubCat string }{
		{"T8_MOUNT_MAMMOTH", "Mamut de Comando (T8)", "Montura Rara"},
		{"UNIQUE_MOUNT_BEAR_KEEPER", "Oso Grizzly", "Montura Rara"},
		{"UNIQUE_MOUNT_GIANT_HORSE", "Caballo Gigante", "Montura Rara"},
		{"UNIQUE_MOUNT_MORGANA_RAVEN", "Cuervo de Morgana", "Montura Rara"},
		{"UNIQUE_MOUNT_RESTING_SALAMANDER", "Salamandra de Pantano", "Montura Rara"},
		{"UNIQUE_MOUNT_BLACK_PANTHER", "Pantera Negra", "Montura Rara"},
		{"T8_BAG_INSIGHT", "Cartera de la Perspicacia (T8)", "Accesorio"},
	}
	for _, l := range luxuryItems {
		items[l.ID] = ItemMetadata{ID: l.ID, Name: l.Name, Tier: 8, Category: "Skin", SubCategory: l.SubCat}
	}

	data, _ := json.MarshalIndent(items, "", "  ")
	os.WriteFile("backend/metadata/metadata.json", data, 0644)
	fmt.Printf("Generated %d items.\n", len(items))

	// Generate a simplified base mapping (ID -> {Cat, SubCat, Name})
	baseMapping := make(map[string]interface{})
	for id, it := range items {
		if !strings.Contains(id, "@") {
			baseMapping[id] = map[string]string{
				"category":     it.Category,
				"sub_category": it.SubCategory,
				"name":         it.Name,
			}
		}
	}
	baseData, _ := json.MarshalIndent(baseMapping, "", "  ")
	os.WriteFile("backend/metadata/base_items.json", baseData, 0644)
}
