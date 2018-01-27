package cgminer

import (
	"testing"
	"strings"
)

func TestProcessChipStat(t *testing.T) {
	json := `,"SUMMARY":[{"STATUS":[{"STATUS":"S","When":1516694992,"Code":11,"Msg":"Summary","Description":"cpuminer 2.3.2"}]{"0_accept":116,"1_accept":108,"2_accept":81,"3_accept":120,"4_accept":120,"5_accept":143,"6_accept":117,"7_accept":0,"8_accept":142,"9_accept":112,"10_accept":105,"11_accept":133,"12_accept":132,"13_accept":54,"14_accept":131,"15_accept":130,"16_accept":122,"17_accept":9,"18_accept":128,"19_accept":122,"20_accept":123,"21_accept":95,"22_accept":102,"23_accept":134,"24_accept":108,"25_accept":67,"26_accept":101,"27_accept":43,"28_accept":37,"29_accept":77,"30_accept":64,"31_accept":127,"32_accept":116,"33_accept":0,"34_accept":0,"35_accept":0,"36_accept":131,"37_accept":105,"38_accept":123,"39_accept":111,"40_accept":128,"41_accept":121,"42_accept":104,"43_accept":7,"44_accept":138,"45_accept":110,"46_accept":41,"47_accept":153,"48_accept":69,"49_accept":84,"50_accept":84,"51_accept":0,"52_accept":59,"53_accept":130,"54_accept":90,"55_accept":108,"56_accept":92,"57_accept":86,"58_accept":122,"59_accept":101,"60_accept":39,"61_accept":41,"62_accept":0,"63_accept":6,"64_accept":8,"65_accept":11,"66_accept":22,"67_accept":44,"68_accept":3,"69_accept":56,"70_accept":0,"71_accept":1,"72_accept":130,"73_accept":101,"74_accept":146,"75_accept":113,"76_accept":105,"77_accept":38,"78_accept":111,"79_accept":12,"80_accept":116,"81_accept":135,"82_accept":117,"83_accept":16,"84_accept":106,"85_accept":117,"86_accept":121,"87_accept":68,"88_accept":133,"89_accept":24,"90_accept":85,"91_accept":133,"92_accept":142,"93_accept":82,"94_accept":5,"95_accept":134,"96_accept":118,"97_accept":61,"98_accept":8,"99_accept":11,"100_accept":20,"101_accept":108,"102_accept":71,"103_accept":137,"104_accept":11,"105_accept":0,"106_accept":4,"107_accept":0,"108_accept":105,"109_accept":104,"110_accept":48,"111_accept":133,"112_accept":72,"113_accept":145,"114_accept":127,"115_accept":0,"116_accept":138,"117_accept":124,"118_accept":31,"119_accept":0,"120_accept":117,"121_accept":102,"122_accept":116,"123_accept":20,"124_accept":114,"125_accept":1,"126_accept":111,"127_accept":7,"128_accept":95,"129_accept":42,"130_accept":103,"131_accept":127,"132_accept":134,"133_accept":120,"134_accept":67,"135_accept":113,"136_accept":94,"137_accept":53,"138_accept":84,"139_accept":111,"140_accept":68,"141_accept":0,"142_accept":1,"143_accept":0,"0_reject":20,"1_reject":26,"2_reject":32,"3_reject":7,"4_reject":5,"5_reject":3,"6_reject":14,"7_reject":1,"8_reject":6,"9_reject":13,"10_reject":6,"11_reject":4,"12_reject":3,"13_reject":68,"14_reject":6,"15_reject":7,"16_reject":11,"17_reject":2,"18_reject":6,"19_reject":4,"20_reject":13,"21_reject":20,"22_reject":7,"23_reject":4,"24_reject":21,"25_reject":68,"26_reject":26,"27_reject":89,"28_reject":83,"29_reject":30,"30_reject":54,"31_reject":2,"32_reject":16,"33_reject":0,"34_reject":0,"35_reject":0,"36_reject":33,"37_reject":20,"38_reject":6,"39_reject":9,"40_reject":7,"41_reject":20,"42_reject":10,"43_reject":2,"44_reject":4,"45_reject":11,"46_reject":4,"47_reject":4,"48_reject":57,"49_reject":42,"50_reject":42,"51_reject":1,"52_reject":37,"53_reject":13,"54_reject":13,"55_reject":38,"56_reject":36,"57_reject":55,"58_reject":10,"59_reject":6,"60_reject":7,"61_reject":13,"62_reject":5,"63_reject":61,"64_reject":62,"65_reject":32,"66_reject":29,"67_reject":10,"68_reject":3,"69_reject":6,"70_reject":3,"71_reject":1,"72_reject":30,"73_reject":32,"74_reject":27,"75_reject":9,"76_reject":52,"77_reject":8,"78_reject":10,"79_reject":1,"80_reject":3,"81_reject":7,"82_reject":31,"83_reject":18,"84_reject":7,"85_reject":12,"86_reject":8,"87_reject":33,"88_reject":7,"89_reject":12,"90_reject":18,"91_reject":5,"92_reject":9,"93_reject":47,"94_reject":10,"95_reject":12,"96_reject":2,"97_reject":60,"98_reject":17,"99_reject":52,"100_reject":75,"101_reject":27,"102_reject":59,"103_reject":10,"104_reject":17,"105_reject":2,"106_reject":1,"107_reject":0,"108_reject":13,"109_reject":7,"110_reject":11,"111_reject":3,"112_reject":56,"113_reject":6,"114_reject":8,"115_reject":0,"116_reject":4,"117_reject":6,"118_reject":2,"119_reject":1,"120_reject":39,"121_reject":56,"122_reject":32,"123_reject":4,"124_reject":25,"125_reject":3,"126_reject":13,"127_reject":8,"128_reject":30,"129_reject":71,"130_reject":7,"131_reject":8,"132_reject":10,"133_reject":19,"134_reject":60,"135_reject":36,"136_reject":46,"137_reject":71,"138_reject":50,"139_reject":27,"140_reject":38,"141_reject":1,"142_reject":3,"143_reject":0}],"id":1}`
	actualResult, err := processChipStat(json)
	if err != nil {
		t.Error(err)
	}

	csLen := len(actualResult.ChipStat)
	if csLen != 288 {
		t.Fatal("Length is not 288")
	}
}

// Fake test for an issue I found. Setting this up for testing in golang is lame.
func TestReadVersionFile(t *testing.T) {
	version := "ltcminer-version-fan-5w-nocheckpool\n"
	originalLength := len(version)

	result := strings.TrimSpace(version)

	if strings.Contains(result, "\n") || originalLength == len(result) {
		t.Fatal("Newline wasn't stripped!")
	}
}
