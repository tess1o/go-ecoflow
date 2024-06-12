package main

type DeviceQuotaResponse struct {
	Code            string `json:"code"`
	Message         string `json:"message"`
	Data            DeviceQuota
	EagleEyeTraceID string `json:"eagleEyeTraceId"`
	Tid             string `json:"tid"`
}

type DeviceQuota struct {
	Mppt         MpptProperties
	BmsEmsStatus BmsEmsStatusProperties
	BmsBmsStatus BmsBmsStatusProperties
	Inv          InvProperties
	Pd           PdProperties
	BmsSlave     BmsSlaveProperties
}

type MpptProperties struct {
	FaultCode     int    `json:"mppt.faultCode"`
	InAmp         int    `json:"mppt.inAmp"`
	Res           [5]int `json:"mppt.res"`
	Dcdc12vAmp    int    `json:"mppt.dcdc12vAmp"`
	CarStandbyMin int    `json:"mppt.carStandbyMin"`
	CarOutAmp     int    `json:"mppt.carOutAmp"`
	ChgType       int    `json:"mppt.chgType"`
	CfgChgType    int    `json:"mppt.cfgChgType"`
	DcChgCurrent  int    `json:"mppt.dcChgCurrent"`
	CarTemp       int    `json:"mppt.carTemp"`
	OutWatts      int    `json:"mppt.outWatts"`
	BeepState     int    `json:"mppt.beepState"`
	SwVer         int    `json:"mppt.swVer"`
	InVol         int    `json:"mppt.inVol"`
	OutVol        int    `json:"mppt.outVol"`
	MpptTemp      int    `json:"mppt.mpptTemp"`
	CfgAcEnabled  int    `json:"mppt.cfgAcEnabled"`
	CfgAcOutVol   int    `json:"mppt.cfgAcOutVol"`
	CfgAcOutFreq  int    `json:"mppt.cfgAcOutFreq"`
	CarOutVol     int    `json:"mppt.carOutVol"`
	Dcdc12vWatts  int    `json:"mppt.dcdc12vWatts"`
	Dc24vTemp     int    `json:"mppt.dc24vTemp"`
	Dc24vState    int    `json:"mppt.dc24vState"`
	PowStandbyMin int    `json:"mppt.powStandbyMin"`
	ScrStandbyMin int    `json:"mppt.scrStandbyMin"`
	InWatts       int    `json:"mppt.inWatts"`
	CarOutWatts   int    `json:"mppt.carOutWatts"`
	DischargeType int    `json:"mppt.dischargeType"`
	ChgPauseFlag  int    `json:"mppt.chgPauseFlag"`
	X60ChgType    int    `json:"mppt.x60ChgType"`
}

type BmsEmsStatusProperties struct {
	MaxChargeSoc     int     `json:"bms_emsStatus.maxChargeSoc"`
	BmsIsConnt       [3]int  `json:"bms_emsStatus.bmsIsConnt"`
	DsgCmd           int     `json:"bms_emsStatus.dsgCmd"`
	ChgAmp           int     `json:"bms_emsStatus.chgAmp"`
	ChgState         int     `json:"bms_emsStatus.chgState"`
	OpenBmsIdx       int     `json:"bms_emsStatus.openBmsIdx"`
	FanLevel         int     `json:"bms_emsStatus.fanLevel"`
	ChgVol           int     `json:"bms_emsStatus.chgVol"`
	ChgRemainTime    int     `json:"bms_emsStatus.chgRemainTime"`
	ParaVolMin       int64   `json:"bms_emsStatus.paraVolMin"`
	ParaVolMax       int     `json:"bms_emsStatus.paraVolMax"`
	MinOpenOilEb     int     `json:"bms_emsStatus.minOpenOilEb"`
	MaxAvailNum      int     `json:"bms_emsStatus.maxAvailNum"`
	LcdShowSoc       int     `json:"bms_emsStatus.lcdShowSoc"`
	BmsModel         int     `json:"bms_emsStatus.bmsModel"`
	EmsIsNormalFlag  int     `json:"bms_emsStatus.emsIsNormalFlag"`
	OpenUpsFlag      int     `json:"bms_emsStatus.openUpsFlag"`
	MaxCloseOilEb    int     `json:"bms_emsStatus.maxCloseOilEb"`
	MinDsgSoc        int     `json:"bms_emsStatus.minDsgSoc"`
	ChgCmd           int     `json:"bms_emsStatus.chgCmd"`
	F32LcdShowSoc    float64 `json:"bms_emsStatus.f32LcdShowSoc"`
	MinDsgRemainTime int     `json:"bms_emsStatus.dsgRemainTime"`
}

type BmsBmsStatusProperties struct {
	DesignCap        int     `json:"bms_bmsStatus.designCap"`
	OutputWatts      int     `json:"bms_bmsStatus.outputWatts"`
	Vol              int     `json:"bms_bmsStatus.vol"`
	RemainCap        int     `json:"bms_bmsStatus.remainCap"`
	MinCellVol       int     `json:"bms_bmsStatus.minCellVol"`
	MaxCellVol       int     `json:"bms_bmsStatus.maxCellVol"`
	MaxCellTemp      int     `json:"bms_bmsStatus.maxCellTemp"`
	MinCellTemp      int     `json:"bms_bmsStatus.minCellTemp"`
	MinMosTemp       int     `json:"bms_bmsStatus.minMosTemp"`
	MaxMosTemp       int     `json:"bms_bmsStatus.maxMosTemp"`
	RemainTime       int     `json:"bms_bmsStatus.remainTime"`
	FullCap          int     `json:"bms_bmsStatus.fullCap"`
	Soc              int     `json:"bms_bmsStatus.soc"`
	SoH              int     `json:"bms_bmsStatus.soh"`
	Cycles           int     `json:"bms_bmsStatus.cycles"`
	TagChgAmp        int     `json:"bms_bmsStatus.tagChgAmp"`
	Amp              int     `json:"bms_bmsStatus.amp"`
	Temp             int     `json:"bms_bmsStatus.temp"`
	MinDsgRemainTime int     `json:"bms_bmsStatus.dsgRemainTime"`
	F32ShowSoc       float64 `json:"bms_bmsStatus.f32ShowSoc"`
	BmsFault         int     `json:"bms_bmsStatus.bmsFault"`
	ErrCode          int     `json:"bms_bmsStatus.errCode"`
	Type             int     `json:"bms_bmsStatus.type"`
	BqSysStatReg     int     `json:"bms_bmsStatus.bqSysStatReg"`
	InputWatts       int     `json:"bms_bmsStatus.inputWatts"`
	OpenBmsIdx       int     `json:"bms_bmsStatus.openBmsIdx"`
	CellId           int     `json:"bms_bmsStatus.cellId"`
	SysVer           int     `json:"bms_bmsStatus.sysVer"`
}

type InvProperties struct {
	InvOutFreq    int    `json:"inv.invOutFreq"`
	CfgAcXboost   int    `json:"inv.cfgAcXboost"`
	OutTemp       int    `json:"inv.outTemp"`
	InputWatts    int    `json:"inv.inputWatts"`
	InvOutAmp     int    `json:"inv.invOutAmp"`
	SlowChgWatts  int    `json:"inv.SlowChgWatts"`
	DcInAmp       int    `json:"inv.dcInAmp"`
	ChargerType   int    `json:"inv.chargerType"`
	CfgAcEnabled  int    `json:"inv.cfgAcEnabled"`
	DischargeType int    `json:"inv.dischargeType"`
	InvOutVol     int    `json:"inv.invOutVol"`
	ErrCode       int    `json:"inv.errCode"`
	AcInVol       int    `json:"inv.acInVol"`
	FastChgWatts  int    `json:"inv.FastChgWatts"`
	InvType       int    `json:"inv.invType"`
	AcInFreq      int    `json:"inv.acInFreq"`
	CfgAcWorkMode int    `json:"inv.cfgAcWorkMode"`
	ChgPauseFlag  int    `json:"inv.chgPauseFlag"`
	AcInAmp       int    `json:"inv.acInAmp"`
	CfgAcOutVol   int    `json:"inv.cfgAcOutVol"`
	StandbyMins   int    `json:"inv.standbyMins"`
	FanState      int    `json:"inv.fanState"`
	OutputWatts   int    `json:"inv.outputWatts"`
	Reserved      [8]int `json:"inv.reserved"`
	DcInVol       int    `json:"inv.dcInVol"`
	DcInTemp      int    `json:"inv.dcInTemp"`
	CfgAcOutFreq  int    `json:"inv.cfgAcOutFreq"`
	AcDipSwitch   int    `json:"inv.acDipSwitch"`
	SysVer        int    `json:"inv.sysVer"`
	InvOutWatts   int    `json:"inv.invOutWatts"`
}

type PdProperties struct {
	ExtRj45Port   int     `json:"pd.extRj45Port"`
	Soc           int     `json:"pd.soc"`
	AcAutoOnCfg   int     `json:"pd.acAutoOnCfg"`
	BrightLevel   int     `json:"pd.brightLevel"`
	Typec2Temp    int     `json:"pd.typec2Temp"`
	Typec1Watts   int     `json:"pd.typec1Watts"`
	ChgDsgState   int     `json:"pd.chgDsgState"`
	Ext3p8Port    int     `json:"pd.ext3p8Port"`
	Typec2Watts   int     `json:"pd.typec2Watts"`
	Typec1Temp    int     `json:"pd.typec1Temp"`
	Ext4p8Port    int     `json:"pd.ext4p8Port"`
	InWatts       int     `json:"pd.inWatts"`
	CarWatts      int     `json:"pd.carWatts"`
	AcEnabled     int     `json:"pd.acEnabled"`
	LcdOffSec     int     `json:"pd.lcdOffSec"`
	WifiAutoRcvy  int     `json:"pd.wifiAutoRcvy"`
	Model         int     `json:"pd.model"`
	WifiVer       int     `json:"pd.wifiVer"`
	DcOutState    int     `json:"pd.dcOutState"`
	RemainTime    int     `json:"pd.remainTime"`
	CarState      int     `json:"pd.carState"`
	LcdShowSoc    int     `json:"pd.lcdShowSoc"`
	UsbUsedTime   int     `json:"pd.usbUsedTime"`
	StandbyMin    int     `json:"pd.standbyMin"`
	BeepMode      int     `json:"pd.beepMode"`
	OutputWatts   int     `json:"pd.outputWatts"`
	WifiRssi      int     `json:"pd.wifiRssi"`
	InvUsedTime   int     `json:"pd.invUsedTime"`
	ChargerType   int     `json:"pd.chargerType"`
	DcInUsedTime  int     `json:"pd.dcInUsedTime"`
	CarTemp       int     `json:"pd.carTemp"`
	SysVer        int     `json:"pd.sysVer"`
	OutWatts      int     `json:"pd.outWatts"`
	WattsOutSum   int     `json:"pd.wattsOutSum"`
	UsbqcUsedTime int     `json:"pd.usbqcUsedTime"`
	PvChgPrioSet  int     `json:"pd.pvChgPrioSet"`
	WattsInSum    int     `json:"pd.wattsInSum"`
	Reserved      [2]int  `json:"pd.reserved"`
	ErrCode       int     `json:"pd.errCode"`
	ChgPowerAC    int     `json:"pd.chgPowerAC"`
	ChgPowerDC    int     `json:"pd.chgPowerDC"`
	QcUsb2Watts   int     `json:"pd.qcUsb2Watts"`
	WireWatts     int     `json:"pd.wireWatts"`
	Usb1Watts     int     `json:"pd.usb1Watts"`
	Usb2Watts     int     `json:"pd.usb2Watts"`
	TypecUsedTime int     `json:"pd.typecUsedTime"`
	WireUsedTime  int     `json:"pd.wireUsedTime"`
	IcoBytes      [14]int `json:"pd.icoBytes"`
}

type BmsSlaveProperties struct {
	MaxCellVol       int `json:"bms_slave.maxCellVol"`
	MinCellVol       int `json:"bms_slave.minCellVol"`
	MinMosTemp       int `json:"bms_slave.minMosTemp"`
	RemainTime       int `json:"bms_slave.remainTime"`
	MaxMosTemp       int `json:"bms_slave.maxMosTemp"`
	Num              int `json:"bms_slave.num"`
	SoH              int `json:"bms_slave.soh"`
	Soc              int `json:"bms_slave.soc"`
	MinCellTemp      int `json:"bms_slave.minCellTemp"`
	BmsFault         int `json:"bms_slave.bmsFault"`
	Cycles           int `json:"bms_slave.cycles"`
	RemainCap        int `json:"bms_slave.remainCap"`
	FullCap          int `json:"bms_slave.fullCap"`
	MinDsgRemainTime int `json:"bms_slave.dsgRemainTime"`
	DesignCap        int `json:"bms_slave.designCap"`
	F32ShowSoc       int `json:"bms_slave.f32ShowSoc"`
	InputWatts       int `json:"bms_slave.inputWatts"`
	OutputWatts      int `json:"bms_slave.outputWatts"`
	Vol              int `json:"bms_slave.vol"`
	SysVer           int `json:"bms_slave.sysVer"`
	Temp             int `json:"bms_slave.temp"`
	CellId           int `json:"bms_slave.cellId"`
	ErrCode          int `json:"bms_slave.errCode"`
	OpenBmsIdx       int `json:"bms_slave.openBmsIdx"`
	TagChgAmp        int `json:"bms_slave.tagChgAmp"`
	BqSysStatReg     int `json:"bms_slave.bqSysStatReg"`
	Type             int `json:"bms_slave.type"`
}
