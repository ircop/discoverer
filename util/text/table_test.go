package text

import (
	"testing"
	"fmt"
)

func Test_ParseShVlanMES(t *testing.T) {
	output := `show vlan
Vlan mode: Basic
Created by: D-Default, S-Static, G-GVRP, R-Radius Assigned VLAN, V-Voice VLAN

show vlan
Vlan mode: Basic
Created by: D-Default, S-Static, G-GVRP, R-Radius Assigned VLAN, V-Voice VLAN

Vlan       Name           Tagged Ports      UnTagged Ports      Created by
---- ----------------- ------------------ ------------------ ----------------
 1           -                            gi2/0/1-24,               D        
                                          te2/0/1-4,                         
                                          gi3/0/1-24,                         
                                          te3/0/1-4,                         
                                          gi4/0/1-24,                         
                                          te4/0/1-4,                         
                                          gi5/0/1-24,                         
                                          te5/0/1-4,                         
                                          gi6/0/1-24,                         
                                          te6/0/1-4,                         
                                          gi7/0/1-24,
                                          te7/0/1-4,
                                          gi8/0/1-24,
                                          te8/0/1-4,Po2-16
999          -          te1/0/1,te1/0/4                             S
1108         -          te1/0/1,te1/0/4                             S
1149         -          te1/0/1,te1/0/4       gi1/0/1-24            S
3380         -          te1/0/1,te1/0/4                             S
3382         -          te1/0/1,te1/0/4                             S
3383         -          te1/0/1,te1/0/4       gi1/0/1-24            S
3384         -          te1/0/1,te1/0/4                             S
3385         -          te1/0/1,te1/0/4                             S
3406         -          te1/0/1,te1/0/4                             S
3408         -          te1/0/1,te1/0/4                             S
3438         -          te1/0/1,te1/0/4                             S
3446         -          te1/0/1,te1/0/4                             S
3447         -          te1/0/1,te1/0/4                             S



`
	//rows := ParseTable(output, `^---`, "")
	rows := ParseTable(output, "^--", "")
	if len(rows) != 14 {
		t.Fatal("Row count should be 14")
	}

	if rows[0][0] != "1" {
		t.Fatalf("Row 0 col 0 should be '1', got '%s'", rows[0][0])
	}
	if rows[0][1] != "-" {
		t.Fatalf("Row 0 col 1 should be '1', got '%s'", rows[0][1])
	}
	if rows[0][2] != "" {
		t.Fatalf("Row 0 col 2 should be '1', got '%s'", rows[0][2])
	}

	fmt.Printf("%+v\n", rows[0][3])

	fmt.Printf("-\n")
}

func Test_ParseShIntStatus(t *testing.T) {
	output := `sh int status

Port      Name               Status       Vlan       Duplex  Speed Type
Gi1/1     -= ATC =-          connected    trunk        full   1000 1000BaseLH
Gi1/2     DXS-112.10         connected    trunk        full   1000 1000BaseLH
Gi1/3     DGS-10.170.112.253 connected    trunk        full   1000 1000BaseBX10-D
Gi1/4     po3:DGS-113.154    connected    trunk        full   1000 1000BaseLH
Gi1/5     DXS-10.170.113.132 notconnect   1            full   1000 No Gbic
Gi1/6     DGS 10.170.113.131 connected    trunk        full   1000 1000BaseBX10-D
Gi1/7     DXS-112.10         connected    trunk        full   1000 1000BaseLH
Gi1/8     Po6 <-> Port 9 DGS connected    trunk        full   1000 1000BaseLH
Gi1/9     DXS-10.170.113.132 connected    trunk        full   1000 Unknown GBIC
Gi1/10    DXS-10.170.113.132 connected    trunk        full   1000 Unknown GBIC
Gi1/11    DGS-10.170.112.253 connected    trunk        full   1000 Unknown GBIC
Gi1/12    DGS-10.170.111.48  connected    trunk        full   1000 Unknown GBIC
Gi1/13    DGS-10.170.113.155 connected    trunk        full   1000 1000BaseLH
Gi1/14    DGS-10.170.111.116 connected    trunk        full   1000 Unknown GBIC
Gi1/15    50-let-VLKSM-12    connected    trunk        full   1000 1000BaseBX10-D
Gi1/16    DGS 10.170.113.100 connected    trunk        full   1000 1000BaseLH
Gi1/17    DGS 10.170.111.157 connected    trunk        full   1000 1000BaseLH
Gi1/18    po3:DGS-113.154    connected    trunk        full   1000 1000BaseLH
Gi1/19    Po1:DGS-10.170.111 connected    trunk        full   1000 1000BaseLH
Gi1/20    Po1:DGS-10.170.111 connected    trunk        full   1000 1000BaseLH
Gi1/21    Po4:111.250        connected    trunk        full   1000 1000BaseLH
Gi1/22    Po4:111.250        connected    trunk        full   1000 1000BaseLH
Gi1/23    Po4:111.250        connected    trunk        full   1000 1000BaseLH
Gi1/24    Po4:111.250        connected    trunk        full   1000 1000BaseLH
Gi1/25    Po7:10.170.114.249 connected    trunk        full   1000 1000BaseLH
Gi1/26    Po6 <-> Port 24 DG notconnect   731          full   1000 1000BaseLH
Gi1/27    Po7:10.170.114.249 connected    trunk        full   1000 1000BaseLH
Gi1/28    -= 10.170.111.48 = connected    trunk        full   1000 Unknown GBIC
Te1/29    -= RAM =-          connected    trunk        full    10G 10GBase-LR
Te1/30    -= Shatura =-      connected    trunk        full    10G 10GBase-LR
Po1       DGS-111.114        connected    trunk      a-full a-1000
Po2       DGS-112.60         notconnect   0            auto   auto
Po3                          connected    trunk      a-full a-1000
Po4       DGS-111.250        connected    trunk      a-full a-1000
Po5       DXS-112.10         connected    trunk      a-full a-1000

`

	rows := ParseTable(output, `^Port\s+Name`, "")
	if len(rows) != 35 {
		t.Fatal("Row count should be 35")
	}

	// each row columns count should be 6
	for i := range rows {
		if len(rows[i]) != 7 {
			t.Fatalf("Row %d column count is %d, should be 7", i, len(rows[i]))
		}
	}

	if rows[6][0] != "Gi1/7" {
		t.Fatalf("Row 7 column 1 should be 'Gi1/7', but got '%s'", rows[6][0])
	}
	if rows[27][6] != "Unknown GBIC" {
		t.Fatalf("Row 27 col 6 should be 'Unknown GBIC', but got '%s'", rows[27][6])
	}
	if rows[31][3] != "0" {
		t.Fatalf("Row 31 (%s) col 3 should be '0', but got '%s'", rows[31][0], rows[31][3])
	}
	if rows[32][1] != "" {
		t.Fatalf("Row 34 col 1 should be '', but got '%s'", rows[32][1])
	}

	/*for i := range rows {
		fmt.Printf("'%s|%s|%s|%s|%s|%s|%s'\n", rows[i][0], rows[i][1],rows[i][2],rows[i][3],rows[i][4],rows[i][5],rows[i][6])
	}*/
}

func Test_ParseShLldpNeiMES(t *testing.T) {
	output := `sh lldp neighbors

System capability legend:
B - Bridge; R - Router; W - Wlan Access Point; T - telephone;
D - DOCSIS Cable Device; H - Host; r - Repeater;
TP - Two Ports MAC Relay; S - S-VLAN; C - C-VLAN; O - Other

  Port        Device ID        Port ID       System Name    Capabilities  TTL
--------- ----------------- ------------- ----------------- ------------ -----
te1/0/1   ac:f1:df:a4:ae:00 ac:f1:df:a4:a                        O        106
                            e:1c
te1/0/4   e0:d9:e3:ba:9e:80    te1/0/3                           O        100

`
	rows := ParseTable(output, `^-----`, "")

/*	for i := range rows {
		fmt.Printf("'%s|%s|%s|%s|%s|%s'\n", rows[i][0], rows[i][1],rows[i][2],rows[i][3],rows[i][4],rows[i][5])
	}*/

	if len(rows) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(rows))
	}
	for i := range rows {
		if len(rows[i]) != 6 {
			t.Fatalf("Row %d: expected 6 cols, got %d", i, len(rows[i]))
		}
	}
	if rows[0][2] != "ac:f1:df:a4:ae:1c" {
		t.Fatalf("Splitted column is broken ('%s')", rows[0][2])
	}
}
