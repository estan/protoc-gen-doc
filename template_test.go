package gendoc_test

import (
	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
	"github.com/stretchr/testify/suite"

	"testing"

	"github.com/pseudomuto/protoc-gen-doc"
)

var (
	template    *gendoc.Template
	bookingFile *gendoc.File
	vehicleFile *gendoc.File
)

type TemplateTest struct {
	suite.Suite
}

func TestTemplate(t *testing.T) {
	suite.Run(t, new(TemplateTest))
}

func (assert *TemplateTest) SetupSuite() {
	set, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	assert.NoError(err)

	req := utils.CreateGenRequest(set, "Booking.proto", "Vehicle.proto")
	result := protokit.ParseCodeGenRequest(req)
	template = gendoc.NewTemplate(result)
	bookingFile = template.Files[0]
	vehicleFile = template.Files[1]
}

func (assert *TemplateTest) TestTemplateProperties() {
	assert.Equal(2, len(template.Files))
}

func (assert *TemplateTest) TestFileProperties() {
	assert.Equal("Booking.proto", bookingFile.Name)
	assert.Equal("Booking related messages.\n\nThis file is really just an example. The data model is completely\nfictional.", bookingFile.Description)
	assert.Equal("com.example", bookingFile.Package)
	assert.True(bookingFile.HasEnums)
	assert.True(bookingFile.HasExtensions)
	assert.True(bookingFile.HasMessages)
	assert.True(bookingFile.HasServices)
}

func (assert *TemplateTest) TestFileEnumProperties() {
	enum := findEnum("BookingStatus.StatusCode", bookingFile)
	assert.Equal("StatusCode", enum.Name)
	assert.Equal("BookingStatus.StatusCode", enum.LongName)
	assert.Equal("com.example.BookingStatus.StatusCode", enum.FullName)
	assert.Equal("A flag for the status result.", enum.Description)
	assert.Equal(2, len(enum.Values))

	expectedValues := []*gendoc.EnumValue{
		{Name: "OK", Number: "200", Description: "OK result."},
		{Name: "BAD_REQUEST", Number: "400", Description: "BAD result."},
	}

	for idx, value := range enum.Values {
		assert.Equal(expectedValues[idx], value)
	}
}

func (assert *TemplateTest) TestFileExtensionProperties() {
	ext := findExtension("BookingStatus.country", bookingFile)
	assert.Equal("country", ext.Name)
	assert.Equal("BookingStatus.country", ext.LongName)
	assert.Equal("com.example.BookingStatus.country", ext.FullName)
	assert.Equal("The country the booking occurred in.", ext.Description)
	assert.Equal("optional", ext.Label)
	assert.Equal("string", ext.Type)
	assert.Equal("string", ext.LongType)
	assert.Equal("string", ext.FullType)
	assert.Equal(100, ext.Number)
	assert.Equal("china", ext.DefaultValue)
	assert.Equal("BookingStatus", ext.ContainingType)
	assert.Equal("BookingStatus", ext.ContainingLongType)
	assert.Equal("com.example.BookingStatus", ext.ContainingFullType)
}

func (assert *TemplateTest) TestMessageProperties() {
	msg := findMessage("Vehicle", vehicleFile)
	assert.Equal("Vehicle", msg.Name)
	assert.Equal("Vehicle", msg.LongName)
	assert.Equal("com.example.Vehicle", msg.FullName)
	assert.Equal("Represents a vehicle that can be hired.", msg.Description)
	assert.False(msg.HasExtensions)
	assert.True(msg.HasFields)
}

func (assert *TemplateTest) TestNestedMessageProperties() {
	msg := findMessage("Vehicle.Category", vehicleFile)
	assert.Equal("Category", msg.Name)
	assert.Equal("Vehicle.Category", msg.LongName)
	assert.Equal("com.example.Vehicle.Category", msg.FullName)
	assert.Equal("Represents a vehicle category. E.g. \"Sedan\" or \"Truck\".", msg.Description)
	assert.False(msg.HasExtensions)
	assert.True(msg.HasFields)
}

func (assert *TemplateTest) TestMessageExtensionProperties() {
	msg := findMessage("Booking", bookingFile)
	assert.Equal(1, len(msg.Extensions))

	ext := msg.Extensions[0]
	assert.Equal("optional_field_1", ext.Name)
	assert.Equal("BookingStatus.optional_field_1", ext.LongName)
	assert.Equal("com.example.BookingStatus.optional_field_1", ext.FullName)
	assert.Equal("An optional field to be used however you please.", ext.Description)
	assert.Equal("optional", ext.Label)
	assert.Equal("string", ext.Type)
	assert.Equal("string", ext.LongType)
	assert.Equal("string", ext.FullType)
	assert.Equal(101, ext.Number)
	assert.Equal("", ext.DefaultValue)
	assert.Equal("BookingStatus", ext.ContainingType)
	assert.Equal("BookingStatus", ext.ContainingLongType)
	assert.Equal("com.example.BookingStatus", ext.ContainingFullType)
	assert.Equal("Booking", ext.ScopeType)
	assert.Equal("Booking", ext.ScopeLongType)
	assert.Equal("com.example.Booking", ext.ScopeFullType)
}

func (assert *TemplateTest) TestFieldProperties() {
	msg := findMessage("BookingStatus", bookingFile)

	field := findField("id", msg)
	assert.Equal("id", field.Name)
	assert.Equal("Unique booking status ID.", field.Description)
	assert.Equal("required", field.Label)
	assert.Equal("int32", field.Type)
	assert.Equal("int32", field.LongType)
	assert.Equal("int32", field.FullType)
	assert.Equal("", field.DefaultValue)
	assert.False(field.IsOneof)

	field = findField("status_code", msg)
	assert.Equal("status_code", field.Name)
	assert.Equal("The status of this status?", field.Description)
	assert.Equal("optional", field.Label)
	assert.Equal("StatusCode", field.Type)
	assert.Equal("BookingStatus.StatusCode", field.LongType)
	assert.Equal("com.example.BookingStatus.StatusCode", field.FullType)
	assert.Equal("", field.DefaultValue)
	assert.False(field.IsOneof)

	field = findField("category", findMessage("Vehicle", vehicleFile))
	assert.Equal("category", field.Name)
	assert.Equal("Vehicle category.", field.Description)
	assert.Equal("", field.Label) // proto3, neither required, nor optional are valid
	assert.Equal("Category", field.Type)
	assert.Equal("Vehicle.Category", field.LongType)
	assert.Equal("com.example.Vehicle.Category", field.FullType)
	assert.Equal("", field.DefaultValue)
	assert.False(field.IsOneof)

	field = findField("properties", findMessage("Vehicle", vehicleFile))
	assert.Equal("properties", field.Name)
	assert.Equal("repeated", field.Label)
	assert.Equal("PropertiesEntry", field.Type)
	assert.Equal("Vehicle.PropertiesEntry", field.LongType)
	assert.Equal("com.example.Vehicle.PropertiesEntry", field.FullType)
	assert.Equal("", field.DefaultValue)
	assert.True(field.IsMap)
	assert.False(field.IsOneof)

	field = findField("rates", findMessage("Vehicle", vehicleFile))
	assert.Equal("rates", field.Name)
	assert.Equal("repeated", field.Label)
	assert.Equal("sint32", field.Type)
	assert.Equal("sint32", field.LongType)
	assert.Equal("sint32", field.FullType)
	assert.False(field.IsMap)
	assert.False(field.IsOneof)

	field = findField("kilometers", findMessage("Vehicle", vehicleFile))
	assert.Equal("kilometers", field.Name)
	assert.Equal("", field.Label)
	assert.Equal("int32", field.Type)
	assert.Equal("int32", field.LongType)
	assert.Equal("int32", field.FullType)
	assert.False(field.IsMap)
	assert.True(field.IsOneof)
	assert.Equal("travel", field.OneofDecl)

	field = findField("human_name", findMessage("Vehicle", vehicleFile))
	assert.Equal("human_name", field.Name)
	assert.Equal("", field.Label)
	assert.Equal("string", field.Type)
	assert.Equal("string", field.LongType)
	assert.Equal("string", field.FullType)
	assert.False(field.IsMap)
	assert.True(field.IsOneof)
	assert.Equal("drivers", field.OneofDecl)
}

func (assert *TemplateTest) TestServiceProperties() {
	service := findService("VehicleService", vehicleFile)
	assert.Equal("VehicleService", service.Name)
	assert.Equal("VehicleService", service.LongName)
	assert.Equal("com.example.VehicleService", service.FullName)
	assert.Equal("The vehicle service.\n\nManages vehicles and such...", service.Description)
	assert.Equal(3, len(service.Methods))
}

func (assert *TemplateTest) TestServiceMethodProperties() {
	service := findService("VehicleService", vehicleFile)

	method := findServiceMethod("AddModels", service)
	assert.Equal("AddModels", method.Name)
	assert.Equal("creates models", method.Description)
	assert.Equal("Model", method.RequestType)
	assert.Equal("Model", method.RequestLongType)
	assert.Equal("com.example.Model", method.RequestFullType)
	assert.Equal(true, method.RequestStreaming)
	assert.Equal("Model", method.ResponseType)
	assert.Equal("Model", method.ResponseLongType)
	assert.Equal("com.example.Model", method.ResponseFullType)
	assert.Equal(true, method.ResponseStreaming)

	method = findServiceMethod("GetVehicle", service)
	assert.Equal("GetVehicle", method.Name)
	assert.Equal("Looks up a vehicle by id.", method.Description)
	assert.Equal("FindVehicleById", method.RequestType)
	assert.Equal("FindVehicleById", method.RequestLongType)
	assert.Equal("com.example.FindVehicleById", method.RequestFullType)
	assert.Equal(false, method.RequestStreaming)
	assert.Equal("Vehicle", method.ResponseType)
	assert.Equal("Vehicle", method.ResponseLongType)
	assert.Equal("com.example.Vehicle", method.ResponseFullType)
	assert.Equal(false, method.ResponseStreaming)
}

func (assert *TemplateTest) TestExcludedComments() {
	message := findMessage("ExcludedMessage", vehicleFile)
	assert.Empty(message.Description)
	assert.Empty(findField("name", message).Description)
	assert.Empty(findField("value", message).Description)

	// just checking that it doesn't exclude everything
	assert.Equal("the id of this message.", findField("id", message).Description)
}

func findService(name string, f *gendoc.File) *gendoc.Service {
	for _, s := range f.Services {
		if s.Name == name {
			return s
		}
	}

	return nil
}

func findServiceMethod(name string, s *gendoc.Service) *gendoc.ServiceMethod {
	for _, m := range s.Methods {
		if m.Name == name {
			return m
		}
	}

	return nil
}

func findEnum(name string, f *gendoc.File) *gendoc.Enum {
	for _, enum := range f.Enums {
		if enum.LongName == name {
			return enum
		}
	}

	return nil
}

func findExtension(name string, f *gendoc.File) *gendoc.FileExtension {
	for _, ext := range f.Extensions {
		if ext.LongName == name {
			return ext
		}
	}

	return nil
}

func findMessage(name string, f *gendoc.File) *gendoc.Message {
	for _, m := range f.Messages {
		if m.LongName == name {
			return m
		}
	}

	return nil
}

func findField(name string, m *gendoc.Message) *gendoc.MessageField {
	for _, f := range m.Fields {
		if f.Name == name {
			return f
		}
	}

	return nil
}
