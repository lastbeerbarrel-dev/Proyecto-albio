import json

items = {}

def add_item(id, name, tier, category, base_ip=0, recipe=None, prices=None):
    items[id] = {
        "id": id,
        "name": name,
        "tier": tier,
        "category": category,
        "base_ip": base_ip,
        "reference_price": 0
    }
    if recipe: items[id]["recipe"] = recipe
    if prices: items[id]["prices"] = prices

resource_info = {"METALBAR": "Lingote", "LEATHER": "Cuero", "CLOTH": "Tela", "PLANKS": "Tablas"}
for t in range(4, 9):
    for r_id, r_name in resource_info.items():
        add_item(f"T{t}_{r_id}", f"{r_name} (T{t})", t, "Resource")

weapon_families = [
    ("SWORD", "Espada", 16, 8, "METALBAR", "LEATHER", ["MAIN_SWORD", "2H_CLAYMORE", "2H_DUALSWORD", "MAIN_SWORD_KEEPER", "2H_SCIMITAR_MORGANA"]),
    ("AXE", "Hacha", 12, 12, "METALBAR", "PLANKS", ["MAIN_AXE", "2H_AXE", "2H_HALBERD", "2H_BEARPAWS", "2H_SCYTHE_HELL"]),
    ("MACE", "Maza", 16, 8, "METALBAR", "CLOTH", ["MAIN_MACE", "2H_MACE", "2H_FLAIL", "MAIN_ROCK_AVALON"]),
    ("DAGGER", "Daga", 12, 12, "LEATHER", "METALBAR", ["MAIN_DAGGER", "2H_DAGGERPAIR", "2H_CLAWSSPELL", "MAIN_DAGGER_KEEPER", "2H_DAGGER_HELL"]),
    ("SPEAR", "Lanza", 16, 8, "PLANKS", "METALBAR", ["MAIN_SPEAR", "2H_SPEAR", "2H_GLAIVE", "MAIN_SPEAR_KEEPER", "2H_HARPOON"]),
    ("BOW", "Arco", 32, 0, "PLANKS", "", ["MAIN_BOW", "2H_WARBOW", "2H_LONGBOW", "2H_BOW_KEEPER", "2H_BOW_HELL"]),
    ("CROSSBOW", "Ballesta", 20, 12, "PLANKS", "METALBAR", ["MAIN_CROSSBOW", "2H_CROSSBOW", "2H_CROSSBOW_LARGE", "MAIN_CROSSBOW_KEEPER"]),
    ("FIRESTAFF", "Bastón de Fuego", 16, 8, "PLANKS", "CLOTH", ["MAIN_FIRESTAFF", "2H_FIRESTAFF", "2H_FIRESTAFF_HELL", "2H_FIRESTAFF_AVALON"]),
    ("FROSTSTAFF", "Bastón de Hielo", 16, 8, "PLANKS", "CLOTH", ["MAIN_FROSTSTAFF", "2H_FROSTSTAFF", "2H_FROSTSTAFF_HELL", "2H_FROSTSTAFF_AVALON"]),
    ("HOLYSTAFF", "Bastón Sagrado", 16, 8, "PLANKS", "CLOTH", ["MAIN_HOLYSTAFF", "2H_HOLYSTAFF", "2H_HOLYSTAFF_HELL", "2H_HOLYSTAFF_AVALON"]),
    ("NATURESTAFF", "Bastón Natural", 16, 8, "PLANKS", "CLOTH", ["MAIN_NATURESTAFF", "2H_NATURESTAFF", "2H_NATURESTAFF_HELL", "2H_NATURESTAFF_AVALON"])
]

for t in range(4, 9):
    for f_id, f_name, c1, c2, r1, r2, variants in weapon_families:
        for v in variants:
            recipe = [{"resource_id": f"T{t}_{r1}", "count": c1}]
            if r2: recipe.append({"resource_id": f"T{t}_{r2}", "count": c2})
            add_item(f"T{t}_{v}", f"{v.replace('_',' ').title()} (T{t})", t, "Weapon", base_ip=300+t*100, recipe=recipe)

armor_families = [
    ("PLATE_SET1", "Soldado", "METALBAR"), ("PLATE_SET2", "Caballero", "METALBAR"), ("PLATE_SET3", "Guardián", "METALBAR"),
    ("LEATHER_SET1", "Mercenario", "LEATHER"), ("LEATHER_SET2", "Cazador", "LEATHER"), ("LEATHER_SET3", "Asesino", "LEATHER"),
    ("CLOTH_SET1", "Erudito", "CLOTH"), ("CLOTH_SET2", "Clérigo", "CLOTH"), ("CLOTH_SET3", "Mago", "CLOTH")
]
for t in range(4, 9):
    for s_id, s_name, res in armor_families:
        for p in ["HEAD", "ARMOR", "SHOES"]:
            recipe = [{"resource_id": f"T{t}_{res}", "count": 8 if p!="ARMOR" else 16}]
            add_item(f"T{t}_{p}_{s_id}", f"{p.title()} de {s_name} (T{t})", t, "Armor", base_ip=300+t*100, recipe=recipe)

offhands = [("OFF_SHIELD", "Escudo", "METALBAR", "PLANKS"), ("OFF_TORCH", "Antorcha", "PLANKS", "CLOTH"), ("OFF_BOOK", "Libro de hechizos", "CLOTH", "")]
for t in range(4, 9):
    for o_id, o_name, r1, r2 in offhands:
        recipe = [{"resource_id": f"T{t}_{r1}", "count": 4}]
        if r2: recipe.append({"resource_id": f"T{t}_{r2}", "count": 4})
        add_item(f"T{t}_{o_id}", f"{o_name} (T{t})", t, "Offhand", base_ip=300+t*100, recipe=recipe)

for t in range(4, 9):
    add_item(f"T{t}_BAG", f"Bolso (T{t})", t, "Accessory", recipe=[{"resource_id": f"T{t}_CLOTH", "count": 8}, {"resource_id": f"T{t}_LEATHER", "count": 8}])
    add_item(f"T{t}_CAPE", f"Capa (T{t})", t, "Accessory", recipe=[{"resource_id": f"T{t}_CLOTH", "count": 4}])
    for m in ["HORSE", "OX"]: add_item(f"T{t}_MOUNT_{m}", f"{m.title()} (T{t})", t, "Mount")

print(json.dumps(items, indent=2, ensure_ascii=False))
