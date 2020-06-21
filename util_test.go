package movielight

import (
	"log"
	"testing"
)

func TestTranslatename(t *testing.T) {
	files := map[string]string{
		"Rocco_und_seine_Brueder_16.12.12_23-15_arte_175_TVOON_DE.mpg.HD.mov":                                                                                                     "Rocco und seine Brueder",
		"Wir_sind_die_Rosinskis_16.11.04_20-15_ard_90_TVOON_DE.mpg.HD.avi":                                                                                                        "Wir sind die Rosinskis",
		"Grand_Budapest_Hotel_17.02.18_20-15_orf1_90_TVOON_DE.mpg.HD.mov":                                                                                                         "Grand Budapest Hotel",
		"Kino-Mord_im_Fahrpreis_inbegriffen-073010-000-A_EQ_0_VA-STA_02772052_MP4-1500_AMM-Tvguide.mp4":                                                                           "Mord im Fahrpreis inbegriffen",
		"Kino-Brian_De_Palma-072417-000-A_EQ_0_VA-STA_02961927_MP4-1500_AMM-PTWEB_170324234631028.mp4":                                                                            "Brian De Palma",
		"Fernsehserie-Stadt_ohne_Namen_(5_6)-054813-005-A_EQ_1_VA-STA_02188101_MP4-1500_AMM-Tvguide.mp4":                                                                          "Stadt ohne Namen (5 6)",
		"Kino_-_Filme-Die_Reise_des_Personalmanagers-041436-000-A_EQ_0_VA-STA_03109597_MP4-1500_AMM-PTWEB_mSdcfxkYn.mp4":                                                          "Die Reise des Personalmanagers",
		"Der_Fernsehfilm_der_Woche-Verräter_-_Tod_am_Meer_-_Nach_dem_Roman_\"Innere_Sicherheit\"_von_Christa_Bernuth-170821_2015_sendung_fdw_2328k_p35v13":                       "Verräter",
		"Das_kleine_Fernsehspiel-Wir_sind_jung._Wir_sind_stark._-_Zum_25._Jahrestag_der_Pogrome_in_Rostock-Lichtenhagen-170817_2215_sendung_dks_2328k_p35v13":                     "Wir sind jung. Wir sind stark.",
		"Taxi_für_eine_Leiche-Taxi_für_eine_Leiche-20171014_2335_sd_02_TAXI-FUER-EINE-_Taxi-fuer-eine-__13949377__o__3440057605__s14153449_9__ORF2HD_23384512P_01010318P_Q6A.mp4": "Taxi für eine Leiche",
		"Filme-Houston-170324105220_houstonneu_118906_webl_ard.mp4":                                                                                                               "Houston",
		"Hard_Sun-Hard_Sun_(3)-180423_2215_sendung_hsn_a1a2_2328k_p35v13.mp4":                                                                                                     "Hard Sun 3",
		"Film_im_rbb-Der_König_von_Berlin-6f535ddd-d218-4488-b831-8eb362b22761_1800k.mp4":                                                                                        "Der König von Berlin",
		"Freistatt.mp4": "Freistatt",
		"Das_alte_Gewehr_-_Spielfilm,_Frankreich_BRD_1975.mp4": "Das alte Gewehr",
	}
	regexes := map[string]string{}
	for key, value := range files {
		log.Println("Key:", key, "Value:", value)
		tName := Translatename(key, regexes)
		if tName != value {
			t.Errorf("translation not matching: result:%s expected:%s", tName, value)
		}
	}
}
