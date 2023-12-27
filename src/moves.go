package main

type Moves struct {
	name        string
	type1       string
	base_power  uint32
	accuracy    uint //100> means it cant miss
	category    string
	priority    int32
	pp          uint32
	crit        float64
	description string
	effect      string
	effect_rate uint32
	target      string //"Selected Target", "Self", "Adjacent Ally", "All Adjacent Foes", "All Adjacent"
	/*IMPORTANT: during data scrape it will instead say "All Adjacent Pokemon" but
	this should be adjusted during data cleaning before being written to the database*/

	//Specific stuff
	contact    bool //All queried as int values but can be represented as bool because they are all 1 or 0
	sound      bool
	punch      bool
	bite       bool
	can_snatch bool
	slice      bool
	bullet     bool
	wind       bool
	powder     bool
	can_metro  bool
	gravity    bool
	defrosts   bool
	bounceable bool
	can_block  bool //Can be blocked by protect or detect
	can_mirror bool
}
