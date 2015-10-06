package svcreg

import (
	"testing"
)

func Test_AddServiceDefinition(t *testing.T) {
	var svdReg = NewGenericServiceRegistry()
	var id, pref, displayName = "mds", "md", "Display"
	var svd, err = svdReg.CreateServiceDefinition(id, pref, displayName)
	if err != nil {
		t.Fatalf("CreateServiceDefinition(%s, %s, %s) returned an error but should not have.", id, pref, displayName)
	}

	if svd.Identifier != id {
		t.Fatalf("CreateServiceDefinition(%s, %s, %s) returned a new instance of ServiceDefinition where member Identifier was set incorrectly.", id, pref, displayName)
	}

	if svd.URIPrefix != pref {
		t.Fatalf("CreateServiceDefinition(%s, %s, %s) returned a new instance of ServiceDefinition where member URIPrefix was set incorrectly.", id, pref, displayName)
	}

	if svd.DisplayName != displayName {
		t.Fatalf("CreateServiceDefinition(%s, %s, %s) returned a new instance of ServiceDefinition where member DisplayName was set incorrectly.", id, pref, displayName)
	}
}

func Test_AddServiceDefinition_DuplicatePrevetion(t *testing.T) {
	var svdReg = NewGenericServiceRegistry()
	var id, pref, displayName = "mds", "md", "Display"
	svdReg.CreateServiceDefinition(id, pref, displayName)
	var _, err = svdReg.CreateServiceDefinition(id, pref, displayName)

	if err == nil {
		t.Fatalf("Expected an error to be returned by a repeated add of a new definition with similar parameters")
	}
}

func Test_AddServiceVersion(t *testing.T) {
	var svdReg = NewGenericServiceRegistry()
	var id, pref, displayName, version = "mds", "md", "Display", "1.0"
	var svd, err = svdReg.CreateServiceDefinition(id, pref, displayName)
	if err != nil {
		t.Fatalf("CreateServiceDefinition(%s, %s, %s) returned an error but should not have.", id, pref, displayName)
	}

	sv, err := svdReg.CreateServiceVersion(svd.Identifier, version)

	if err != nil {
		t.Fatalf("CreateServiceVersion(%s, %s) returned an error but should not have: %s", id, version, err.Error())
	} else if sv.Version != version {
		t.Fatalf("CreateServiceVersion(%s, %s) returned a new instance of ServiceVersion where member Version was set incorrectly.", id, version)
	} else if sv.ServiceIdentifier != id {
		t.Fatalf("CreateServiceVersion(%s, %s) returned a new instance of ServiceVersion where member ServiceIdentifier was set incorrectly.", id, version)
	}
}

func Test_AddServiceVersion_DuplicatePrevetion(t *testing.T) {
	var svdReg = NewGenericServiceRegistry()
	var id, pref, displayName, version = "mds", "md", "Display", "1.0"
	var svd, err = svdReg.CreateServiceDefinition(id, pref, displayName)
	if err != nil {
		t.Fatalf("CreateServiceDefinition(%s, %s, %s) returned an error but should not have.", id, pref, displayName)
	}

	_, err = svdReg.CreateServiceVersion(svd.Identifier, version)
	_, err = svdReg.CreateServiceVersion(svd.Identifier, version)

	if err == nil {
		t.Fatalf("Expected an error to be returned by a repeated add of a version with similar parameters")
	}
}
