package main

import ()

// Finish later
func GetAbilities() map[string]string {
	abilities := map[string]string{ //Types: onswap, ondamagedeal, ondamagerecieve, onfaint, onturnend, onstatlower, onstatus
		"Adaptibility": "ondamagedeal",
		"Aftermath":    "onfaint",
		"Air Lock":     "onswap",
		"Anger Point":  "ondamagerecieve", //MAY NEED SPECIAL TREATEMENT
		"Anticipation": "onswap",
		"Arena Trap":   "onswap",
		"Bad Dreams":   "onturnend", //MAY NEED SPECIAL TREATEMENT (might be easier to just check during damage calc)
		"Battle Armor": "ondamagedeal",
		"Blaze":        "ondamagedeal",
		"Cacophony":    "ondamagerecieve",
		"Chlorophyll":  "onswap", //MAY NEED SPECIAL TREATEMENT
		"Clear Body":   "onstatlower",
		"Cloud Nine":   "onswap",
		"Color Change": "ondamagerecieve",
		"Compoundeyes": "ondamagedeal",
		"Cute Charm":   "ondamagerecieve",
		"Damp":         "ondamagerecieve",
		"Download":     "onswap",
		"Drizzle":      "onswap",
		"Drought":      "onswap",
		"Dry Skin":     "onturnend",
		"Early Bird":   "ondamagedeal",
		"Effect Spore": "ondamagerecieve",
		"Filter":       "ondamagerecieve",
		"Flame Body":   "ondamagerecieve",
		"Flash Fire":   "ondamagerecieve",
		"Flower Gift":  "onswap", //MAY NEED SPECIAL TREATEMENT
		"Forecast":     "onswap", //Special castform ability
		"Forewarn":     "onswap",
		"Frisk":        "onswap",
		"Gluttony":     "ondamagerecieve",
		"Guts":         "onstatus",
		"Heat Proof":   "ondamagerecieve",
		//"Honey Gather": "?????" Leave blank for now because it has no in battle effect
		"Huge Power":   "onswap", //Same situation as guts
		"Hustle":       "ondamagedeal",
		"Hydration":    "onturnend",
		"Hyper Cutter": "onstatlower",
		"Ice Body":     "onturnend",
		//"Illuminate": "????" No in battle effect
		"Immunity":    "onstatus",
		"Inner Focus": "ondamagerecieve",
		"Insomnia":    "onstatus",
		"Intimidate":  "onswap",
		"Iron Fist":   "ondamagedeal",
		"Keen Eye":    "onstatlower",
		//"Klutz": Special case, on item activation check for klutz
		"Leaf Guard":    "onturnend",
		"Levitate":      "ondamagerecieve",
		"Lightning Rod": "ondamagerecieve",
		"Limber":        "onstatus",
		"Liquid Ooze":   "ondamagerecieve",
		"Magic Guard":   "ondamagerecieve", //Special case!!
		"Magma Armor":   "onstatus",
		"Magnet Pull":   "onswap",
		"Marvel Scale":  "onstatus",
		//"Minus": "???" Special case
		//"Mold Breaker": "???" Special case
		"Motor Drive":  "ondamagerecieve",
		"MultiType":    "onswap",
		"Natural Cure": "onswap", //Special case because its a swap out???
		//"No Guard": "ondamagerecieve" //Special case because it applies to damage dealt and taken
		"Normalize": "ondamagedeal",
		"Oblivious": "onstatus",
		"Overgrow":  "ondamagedeal",
		"Own Tempo": "onstatus",
		//"PickUp": "???" No battle effect
		//"Plus": "???" Special case with minus
		"Poison Heal":  "onturnend",
		"Poison Point": "ondamagedeal",
		"Pressure":     "ondamagerecieve",
		"Pure Power":   "onswap",
		"Quick Feet":   "onstatus",
		"Reckless":     "ondamagedeal",
		"Rivalry":      "ondamagedeal",
		"Rain Dish":    "onturnend",
		"Rock Head":    "ondamagedeal",
		"Rough Skin":   "ondamagerecieve",
		//"Run Away": "???" no in battle effect
		"Sand Stream":  "onswap",
		"Sand Veil":    "ondamagerecieve",
		"Scrappy":      "ondamagedeal",
		"Serene Grace": "ondamagedeal",
		"Shadow Tag":   "onswap",
		"Shed Skin":    "onturnend",
		"Shell Armor":  "ondamagerecieve", //See battle armor
		"Shield Dust":  "ondamagerecieve", //Special treatment
		//"Simple": "???" Not sure atm
		"Skill Link":   "ondamagedeal",
		"Slow Start":   "onturnend", //Special case
		"Sniper":       "ondamagedeal",
		"Snow Cloak":   "ondamagerecieve",
		"Snow Warning": "onswap",
		"Solar Power":  "onswap", //Special case, 2 keys
		"Solid Rock":   "ondamagerecieve",
		"Soundproof":   "ondamagerecieve",
		"Speed Boost":  "onturnend",
		"Static":       "ondamagerecieve",
		//"Stall": "???" Not sure atm
		"Steadfast": "ondamagerecieve",
		//"Stench": "???" No battle uses
		"Sticky Hold":  "ondamagerecieve",
		"Storm Drain":  "ondamagerecieve", //Special case
		"Sturdy":       "ondamagerecieve",
		"Suction Cups": "ondamagerecieve",
		"Super Luck":   "ondamagedeal",
		"Swarm":        "ondamagedeal",
		//"Swift Swim": "???" On swap on weather?
		"Synchronize":  "onstatus",
		"Tangled Feet": "ondamagerecieve",
		"Technician":   "ondamagedeal",
		"Thick Fat":    "ondamagerecieve",
		"Tinted Lens":  "ondamagedeal",
		"Torrent":      "ondamagedeal",
		"Trace":        "onswap",
		"Truant":       "onturnend",
		//"Unaware": "???" Special case
		//"Unburden": "???" Special case
		"Vital Spirit": "onstatus",
		"Volt Absorb":  "ondamagerecieve",
		"Water Absorb": "ondamagerecieve",
		"Water Veil":   "onstatus",
		"White Smoke":  "onstatlower",
		"Wonder Guard": "ondamagerecieve",
		//More are missing such as from newer games (Water bubble, beast boost)
	}
	/*Consider removing ones that no pokemon in the game possess*/
	return abilities
}
