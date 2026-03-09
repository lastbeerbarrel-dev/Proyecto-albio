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

# T4 to T8 base resources
resource_info = {
    "METALBAR": "Lingote de metal",
    "LEATHER": "Cuero",
    "CLOTH": "Tela",
    "PLANKS": "Tablas"
}
for t in range(4, 9):
    for r_id, r_name in resource_info.items():
        add_item(f"T{t}_{r_id}", f"{r_name} (T{t})", t, "Resource")

# Weapons - Multi-tier
weapon_types = [
    ("MAIN_SWORD", "Espada ancha", 16, 8, "METALBAR", "LEATHER"),
    ("2H_CLAYMORE", "Claymore", 20, 12, "METALBAR", "LEATHER"),
    ("2H_DUALSWORD", "Espadas dobles", 20, 12, "METALBAR", "LEATHER"),
    ("MAIN_SWORD_KEEPER", "Hoja de Clarent", 16, 8, "METALBAR", "LEATHER"),
    ("2H_SCIMITAR_MORGANA", "Espada tallada", 20, 12, "METALBAR", "LEATHER"),
    ("MAIN_AXE", "Hacha de batalla", 12, 12, "METALBAR", "PLANKS"),
    ("2H_AXE", "Gran hacha", 16, 16, "METALBAR", "PLANKS"),
    ("2H_HALBERD", "Alabarda", 16, 16, "METALBAR", "PLANKS"),
    ("MAIN_MACE", "Maza", 16, 8, "METALBAR", "CLOTH"),
    ("2H_MACE", "Maza pesada", 20, 12, "METALBAR", "CLOTH"),
    ("MAIN_DAGGER", "Daga", 12, 12, "LEATHER", "METALBAR"),
    ("2H_DAGGERPAIR", "Dagas dobles", 16, 16, "LEATHER", "METALBAR"),
    ("2H_CLAWSSPELL", "Garras", 16, 16, "LEATHER", "METALBAR"),
    ("MAIN_DAGGER_KEEPER", "Sangradora", 12, 12, "LEATHER", "METALBAR"),
    ("MAIN_SPEAR", "Lanza", 16, 8, "PLANKS", "METALBAR"),
    ("2H_SPEAR", "Pica", 20, 12, "PLANKS", "METALBAR"),
    ("2H_GLAIVE", "Glaive", 20, 12, "PLANKS", "METALBAR"),
    ("MAIN_FIRESTAFF", "Bastón de fuego", 16, 8, "PLANKS", "CLOTH"),
    ("2H_FIRESTAFF", "Gran bastón de fuego", 20, 12, "PLANKS", "CLOTH"),
    ("MAIN_FROSTSTAFF", "Bastón de hielo", 16, 8, "PLANKS", "CLOTH"),
    ("2H_FROSTSTAFF", "Gran bastón de hielo", 20, 12, "PLANKS", "CLOTH"),
    ("MAIN_BOW", "Arco", 32, 0, "PLANKS", ""),
    ("2H_WARBOW", "Arco de guerra", 32, 0, "PLANKS", ""),
    ("2H_LONGBOW", "Arco largo", 32, 0, "PLANKS", "")
]

for t in range(4, 9):
    for w_id, w_name, c1, c2, r1, r2 in weapon_types:
        recipe = [{"resource_id": f"T{t}_{r1}", "count": c1}]
        if r2: recipe.append({"resource_id": f"T{t}_{r2}", "count": c2})
        add_item(f"T{t}_{w_id}", f"{w_name} (T{t})", t, "Weapon", base_ip=300+t*100, recipe=recipe)

# Armor Sets - Multi-tier
armor_sets = [
    ("PLATE_SET1", "Soldado", "METALBAR"),
    ("PLATE_SET2", "Caballero", "METALBAR"),
    ("PLATE_SET3", "Guardián", "METALBAR"),
    ("LEATHER_SET1", "Mercenario", "LEATHER"),
    ("LEATHER_SET2", "Cazador", "LEATHER"),
    ("LEATHER_SET3", "Asesino", "LEATHER"),
    ("CLOTH_SET1", "Erudito", "CLOTH"),
    ("CLOTH_SET2", "Clérigo", "CLOTH"),
    ("CLOTH_SET3", "Mago", "CLOTH")
]

for t in range(4, 9):
    for s_id, s_name, res in armor_sets:
        # Head (8)
        add_item(f"T{t}_HEAD_{s_id}", f"Casco de {s_name} (T{t})", t, "Armor", base_ip=300+t*100, recipe=[{"resource_id": f"T{t}_{res}", "count": 8}])
        # Chest (16)
        add_item(f"T{t}_ARMOR_{s_id}", f"Peto de {s_name} (T{t})", t, "Armor", base_ip=300+t*100, recipe=[{"resource_id": f"T{t}_{res}", "count": 16}])
        # Shoes (8)
        add_item(f"T{t}_SHOES_{s_id}", f"Botas de {s_name} (T{t})", t, "Armor", base_ip=300+t*100, recipe=[{"resource_id": f"T{t}_{res}", "count": 8}])

# Offhands
offhands = [
    ("OFF_SHIELD", "Escudo", "METALBAR", "PLANKS"),
    ("OFF_TORCH", "Antorcha", "PLANKS", "CLOTH"),
    ("OFF_BOOK", "Libro de hechizos", "CLOTH", "")
]
for t in range(4, 9):
    for o_id, o_name, r1, r2 in offhands:
        recipe = [{"resource_id": f"T{t}_{r1}", "count": 4}]
        if r2: recipe.append({"resource_id": f"T{t}_{r2}", "count": 4})
        add_item(f"T{t}_{o_id}", f"{o_name} (T{t})", t, "Offhand", base_ip=300+t*100, recipe=recipe)

# Accessories
for t in range(4, 9):
    add_item(f"T{t}_BAG", f"Bolso (T{t})", t, "Accessory", recipe=[{"resource_id": f"T{t}_CLOTH", "count": 8}, {"resource_id": f"T{t}_LEATHER", "count": 8}])
    add_item(f"T{t}_CAPE", f"Capa (T{t})", t, "Accessory", recipe=[{"resource_id": f"T{t}_CLOTH", "count": 4}])

# Mounts
for t in [4,5,6,7,8]:
    add_item(f"T{t}_MOUNT_HORSE", f"Caballo de montar (T{t})", t, "Mount")
    add_item(f"T{t}_MOUNT_OX", f"Buey de transporte (T{t})", t, "Mount")

with open('metadata.json', 'w', encoding='utf-8') as f:
    json.dump(items, f, indent=2, ensure_ascii=False)

print(f"Generated {len(items)} items.")
