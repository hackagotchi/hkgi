
const manifest = {
  items: {
    "bbc_item":         { name: "Rolling Pin",
                          desc: "Makes all your bracti grow more! Probably ethical!" },
    "bbc_egg":          { name: "Bread Egg",
                          usable: true,
                          desc: "Contains hacker and coffee themed goodies, maybe dirt!" },
    "bbc_compressence": { name: "bressence",
                          desc: "EGGREDIENT. More bread than is safe to consume at once." },
    "bbc_essence":      { name: "Bread Essence",
                          desc: "COMPRESSABLE. Baked by a cactus. Tad undercooked." },
    "bbc_seed":         { name: "Bractus Seed",
                          desc: "Doughy, yet somehow still prickly." },
                       
    "hvv_item":         { name: "VINEB0RD",
                          desc: "Makes all your HVVs grow more! CLACK. CL4CK. CLACK." },
    "hvv_egg":          { name: "H4CKER 3GG",
                          usable: true,
                          desc: "Contains bread and coffee themed goodies, maybe dirt!" },
    "hvv_compressence": { name: "hacksprit",
                          desc: "EGGREDIENT. Hacker Spirit compressed with duct tape!" },
    "hvv_essence":      { name: "H4CK3R SP1RIT",
                          desc: "COMPRESSABLE. Makes you wanna make a thing! " },
    "hvv_seed":         { name: "HVV S33D",
                          desc: "Grows into 1337 h4x0r v1n3!" },
                       
    "cyl_item":         { name: "Cyl Wand",
                          desc: "Makes all your cyls grow more! Magic coffee stirrer!" },
    "cyl_egg":          { name: "Cyl Egg",
                          usable: true,
                          desc: "Contains bread and hacker themed goodies, maybe dirt!" },
    "cyl_compressence": { name: "crystcyl",
                          desc: "EGGREDIENT. Hums faintly. Can roast marshmallows!" },
    "cyl_essence":      { name: "Cyl Crystal",
                          desc: "COMPRESSABLE. Glows with aromatic orange energy!" },
    "cyl_seed":         { name: "Cyl Seed",
                          desc: "Do normal coffee beans glow orange in the dark?" },
                       
    "nest_egg":         { name: "Nest Egg",
                          usable: true,
                          desc: "OPEN ME. Some stuff to help you get started :)" },
                       
    "powder_t1":        { name: "Warp Powder",
                          usable: true,
                          desc: "Sparkles, glitters, glows, accelerates plant growth!" },
    "powder_t2":        { name: "Rift Powder",
                          usable: true,
                          desc: "x10 better Warp Powder! Rips open space time." },
    "powder_t3":        { name: "Wormhole Powder",
                          usable: true,
                          desc: "x100 better Warp Powder! May summon time worm." },
                         
    "land_deed":        { name: "Land Deed",
                          usable: true,
                          desc: "Lets you grow more things! Basically dirt paper!" },
   
   /* bagling update */
   "cyl_bag_t1":        { name: "Crystalline Buzzwing Bagling",
                          desc: "Bag. Orange. Small. Flies. Seems a bit seedy." },
   "hvv_bag_t1":        { name: "Spirited Buzzwing Bagling",
                          desc: "Bag. Green. Small. Flies. Seems a bit seedy." },
   "bbc_bag_t1":        { name: "Doughy Buzzwing Bagling",
                          desc: "Bag. Brown. Small. Flies. Seems a bit seedy." },
   "cyl_bag_t2":        { name: "Crystalline Buzzwing Boxlet",
                          desc: "Box. Orange. Small. Flies. Seems seedy." },
   "hvv_bag_t2":        { name: "Spirited Buzzwing Boxlet",
                          desc: "Box. Green. Small. Flies. Seems seedy." },
   "bbc_bag_t2":        { name: "Doughy Buzzwing Boxlet",
                          desc: "Box. Brown. Small. Flies. Seems seedy." },
   "cyl_bag_t3":        { name: "Crystalline Buzzwing Megabox",
                          desc: "Chest. Orange. Winged. Too fat to fly. Quite seedy." },
   "hvv_bag_t3":        { name: "Spirited Buzzwing Megabox",
                          desc: "Chest. Green. Winged. Too fat to fly. Quite seedy." },
   "bbc_bag_t3":        { name: "Doughy Buzzwing Megabox",
                          desc: "Chest. Brown. Winged. Too fat to fly. Quite seedy." },

   "bag_egg_t1":        { name: "Bagling Egg",
                          usable: true,
                          desc: "Boons birthed of a Bagling!" },
   "bag_egg_t2":        { name: "Boxlet Egg",
                          usable: true,
                          desc: "Boons birthed of a Boxlet! Bountiful!" },
   "bag_egg_t3":        { name: "Megabox Egg",
                          usable: true,
                          desc: "The best boons born by a winged vessel!" },
 },
 plant_titles: {
   "dirt": "Dirt",
   "bbc":  "Bractus",
   "cyl":  "Coffea Cyl Plant",
   "hvv":  "H4CK3R V1B3Z V1N3",
 },
 plant_recipes: (() => {
   const trioPlant = (slug, frens) => {
     const compressence = {
       lvl: 0,
       needs: { [slug + "_essence"]: 5 },
       make_item: slug + "_compressence"
     };
     const egg = { lvl: 0, needs: {}, change_plant_to: "dirt", make_item: slug + "_egg" };
     egg.needs[frens[0] + "_compressence"] = 4;
     egg.needs[frens[1] + "_compressence"] = 4;

     const bag_t1 = {
       lvl: 14,
       needs: { [slug + "_compressence"]: 10, [slug + "_bag_t1"]: 3 },
       make_item: {
         one_of: [
           [0.40, "bag_egg_t1"],
           [0.18, frens[0] + "_bag_t2"],
           [0.18, frens[1] + "_bag_t2"],
           [0.12, frens[0] + "_bag_t1"],
           [0.12, frens[1] + "_bag_t1"],
         ]
       }
     };

     const bag_t2 = {
       lvl: 17,
       needs: { [slug + "_compressence"]: 50, [slug + "_bag_t2"]: 2 },
       make_item: {
         one_of: [
           [0.32, "bag_egg_t2"],
           [0.15, frens[0] + "_bag_t3"],
           [0.15, frens[1] + "_bag_t3"],
           [0.19, frens[0] + "_bag_t2"],
           [0.19, frens[1] + "_bag_t2"],
         ]
       }
     };

     const bag_t3 = {
       lvl: 21,
       needs: { [slug + "_compressence"]: 200, [slug + "_bag_t3"]: 1 },
       make_item: {
         one_of: [
           [0.16, "bag_egg_t3"],
           [0.42, frens[0] + "_bag_t3"],
           [0.42, frens[1] + "_bag_t3"],
         ]
       }
     };

     return [compressence, egg, bag_t1, bag_t2, bag_t3];
   };
   return {
     "dirt": [
       { lvl: 0, needs: { "bbc_seed": 1 }, change_plant_to: "bbc" },
       { lvl: 0, needs: { "hvv_seed": 1 }, change_plant_to: "hvv" },
       { lvl: 0, needs: { "cyl_seed": 1 }, change_plant_to: "cyl" },
     ],
     bbc: trioPlant("bbc", [       "cyl", "hvv"]),
     cyl: trioPlant("cyl", ["bbc",        "hvv"]),
     hvv: trioPlant("hvv", ["bbc", "cyl",      ]),
   };
 })(),
};

console.log(manifest)
