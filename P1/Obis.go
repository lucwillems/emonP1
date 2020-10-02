package P1

// variable types used by OBIS
const (
	String    = "string"
	Hex       = "hex"
	Float     = "float"
	MBus      = "mbus"
	Bool      = "bool"
	Integer   = "int"
	Timestamp = "time"
)

type OBISId string

const (
	OBISTypeNil                           = "0-0:0.0.0"
	OBISTypeDateTimestamp                 = "0-0:1.0.0"
	OBISTypeEquipmentIdentifier           = "0-0:96.1.1"
	OBISTypeBEVersionInfo                 = "0-0:96.1.4"
	OBISTypeNumberOfPowerFailures         = "0-0:96.7.21"
	OBISTypeNumberOfLongPowerFailures     = "0-0:96.7.9"
	OBISTypeTextMessage                   = "0-0:96.13.0"
	OBISTypeElectricityTariffIndicator    = "0-0:96.14.0"
	OBISTypeElectricityDelivered          = "1-0:1.7.0"
	OBISTypeElectricityDeliveredTariff1   = "1-0:1.8.1"
	OBISTypeElectricityDeliveredTariff2   = "1-0:1.8.2"
	OBISTypeElectricityGenerated          = "1-0:2.7.0"
	OBISTypeElectricityGeneratedTariff1   = "1-0:2.8.1"
	OBISTypeElectricityGeneratedTariff2   = "1-0:2.8.2"
	OBISTypePowerFailureEventLog          = "1-0:99.97.0"
	OBISTypeInstantaneousPowerDeliveredL1 = "1-0:21.7.0"
	OBISTypeInstantaneousPowerGeneratedL1 = "1-0:22.7.0"
	OBISTypeInstantaneousCurrentL1        = "1-0:31.7.0"
	OBISTypeInstantaneousVoltageL1        = "1-0:32.7.0"
	OBISTypeNumberOfVoltageSagsL1         = "1-0:32.32.0"
	OBISTypeNumberOfVoltageSwellsL1       = "1-0:32.36.0"
	OBISTypeInstantaneousPowerDeliveredL2 = "1-0:41.7.0"
	OBISTypeInstantaneousPowerGeneratedL2 = "1-0:42.7.0"
	OBISTypeInstantaneousCurrentL2        = "1-0:51.7.0"
	OBISTypeInstantaneousVoltageL2        = "1-0:52.7.0"
	OBISTypeNumberOfVoltageSagsL2         = "1-0:52.32.0"
	OBISTypeNumberOfVoltageSwellsL2       = "1-0:52.36.0"
	OBISTypeInstantaneousPowerDeliveredL3 = "1-0:61.7.0"
	OBISTypeInstantaneousPowerGeneratedL3 = "1-0:62.7.0"
	OBISTypeInstantaneousCurrentL3        = "1-0:71.7.0"
	OBISTypeInstantaneousVoltageL3        = "1-0:72.7.0"
	OBISTypeNumberOfVoltageSwellsL3       = "1-0:72.36.0"
	OBISTypeNumberOfVoltageSagsL3         = "1-0:72.32.0"
	OBISTypeVersionInformation            = "1-3:0.2.8"
	//GAS to be verify
	OBISTypeGasEquipmentIdentifier       = "0-1:96.1.0"
	OBISTypeGasEquipmentIdentifierBE     = "0-1:96.1.1"
	OBISTypeGasDeviceType                = "0-1:24.1.0"
	OBISTypeGasTempCorrectedDelivered    = "0-1:24.2.1"
	OBISTypeGasTempNotCorrectedDelivered = "0-1:24.2.3"
	OBISTypeGasValveState                = "0-1:24.4.0"
)

type OType struct {
	Type        string
	Unit        string
	Description string
}

var (
	TypeInfo = map[string]OType{
		//this is internal NIL OBIS Type
		OBISTypeNil: {String, "", "NIL/Invalid type"},
		//common OBIS types
		OBISTypeVersionInformation:            {String, "", "Version Information"},
		OBISTypeDateTimestamp:                 {Timestamp, "", "Date timestamp"},
		OBISTypeEquipmentIdentifier:           {Hex, "", "Equipment Identifier"},
		OBISTypeElectricityDeliveredTariff1:   {Float, "KWh", "Electricity delivered to client (tariff 1)"},
		OBISTypeElectricityDeliveredTariff2:   {Float, "KWh", "Electricity delivered to client (tariff 2)"},
		OBISTypeElectricityGeneratedTariff1:   {Float, "KWh", "Electricity generated by client (tariff 1)"},
		OBISTypeElectricityGeneratedTariff2:   {Float, "KWh", "Electricity generated by client (tariff 2)"},
		OBISTypeElectricityTariffIndicator:    {Integer, "", "Electricity tariff indicator"},
		OBISTypeElectricityDelivered:          {Float, "KW", "Actual electricity delivered"},
		OBISTypeElectricityGenerated:          {Float, "KW", "Actual electricity generated"},
		OBISTypeNumberOfPowerFailures:         {Integer, "", "Number of power failures on any phase"},
		OBISTypeNumberOfLongPowerFailures:     {Integer, "", "Number of long power failures on any phase"},
		OBISTypePowerFailureEventLog:          {String, "", "Event log for long power failures"},
		OBISTypeNumberOfVoltageSagsL1:         {Integer, "", "Number of voltage sags on phase L1"},
		OBISTypeNumberOfVoltageSagsL2:         {Integer, "", "Number of voltage sags on phase L2"},
		OBISTypeNumberOfVoltageSagsL3:         {Integer, "", "Number of voltage sags on phase L3"},
		OBISTypeNumberOfVoltageSwellsL1:       {Integer, "", "Number of voltage swells on phase L1"},
		OBISTypeNumberOfVoltageSwellsL2:       {Integer, "", "Number of voltage swells on phase L2"},
		OBISTypeNumberOfVoltageSwellsL3:       {Integer, "", "Number of voltage swells on phase L3"},
		OBISTypeTextMessage:                   {Hex, "", "Text message"},
		OBISTypeInstantaneousVoltageL1:        {Float, "V", "Instantaneous voltage on phase L1"},
		OBISTypeInstantaneousVoltageL2:        {Float, "V", "Instantaneous voltage on phase L2"},
		OBISTypeInstantaneousVoltageL3:        {Float, "V", "Instantaneous voltage on phase L3"},
		OBISTypeInstantaneousCurrentL1:        {Float, "A", "Instantaneous current on phase L1"},
		OBISTypeInstantaneousCurrentL2:        {Float, "A", "Instantaneous current on phase L2"},
		OBISTypeInstantaneousCurrentL3:        {Float, "A", "Instantaneous current on phase L3"},
		OBISTypeInstantaneousPowerDeliveredL1: {Float, "KW", "Instantaneous active power delivered on phase L1"},
		OBISTypeInstantaneousPowerDeliveredL2: {Float, "KW", "Instantaneous active power delivered on phase L2"},
		OBISTypeInstantaneousPowerDeliveredL3: {Float, "KW", "Instantaneous active power delivered on phase L3"},
		OBISTypeInstantaneousPowerGeneratedL1: {Float, "KW", "Instantaneous active power generated on phase L1"},
		OBISTypeInstantaneousPowerGeneratedL2: {Float, "KW", "Instantaneous active power generated on phase L2"},
		OBISTypeInstantaneousPowerGeneratedL3: {Float, "KW", "Instantaneous active power generated on phase L3"},

		//GAS to be verify
		OBISTypeGasEquipmentIdentifier:       {Hex, "", "Equipment Identifier (NL)"},
		OBISTypeGasEquipmentIdentifierBE:     {Hex, "", "Equipment Identifier (BE)"},
		OBISTypeGasDeviceType:                {Integer, "", "Device Type"},
		OBISTypeGasTempNotCorrectedDelivered: {Float, "m3", "Not temperature corrected volume gas delivered"},
		OBISTypeGasTempCorrectedDelivered:    {MBus, "m3", "Temperature corrected volume gas delivered"},
		OBISTypeGasValveState:                {Integer, "", "Valve state"},
	}
)
