package db

func addOtto(o *Otto) {
	var tmp Otto
	DB.FirstOrCreate(&tmp, o)
}

func fillMockData() {
	addOtto(&Otto{Serial: "0ABCDE", OTPSecret: "JBSWY3DPEHPK3PXP"})
	addOtto(&Otto{Serial: "012345", OTPSecret: "IAWUCNAYW12ANV1K"})
}
