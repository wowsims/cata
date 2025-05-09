package core

import (
	"errors"
	"os"
	"testing"

	"github.com/wowsims/mop/sim/core/proto"
	"google.golang.org/protobuf/encoding/prototext"
)

func readExpectedProtoVersion(fileName string, allowMissingFile bool) (int32, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) && allowMissingFile {
			return 0, nil
		}

		return -1, err
	}

	savedVersionMessage := &proto.ProtoVersion{}

	if err = prototext.Unmarshal(data, savedVersionMessage); err != nil {
		return -1, err
	}

	return savedVersionMessage.SavedVersionNumber, err
}

func TestProtoVersioning(t *testing.T) {
	// // First run the "buf breaking" utility to determine whether any breaking proto changes have been made compared to the
	// // remote master.
	// cmd := exec.Command("npx", "buf", "breaking", "--against", "https://github.com/wowsims/mop.git#branch=master,subdir=proto")
	// cmd.Dir = "../../"
	// out, err := cmd.CombinedOutput()

	// // The utility returns exit code 100 if a breaking change was
	// // identified, so catch only this exit code value.
	// breakingChangeDetected := false

	// if err != nil {
	// 	exiterr, ok := err.(*exec.ExitError)

	// 	if !ok || (exiterr.ExitCode() != 100) {
	// 		t.Fatalf("buf breaking error: %v", err)
	// 	} else {
	// 		breakingChangeDetected = true
	// 	}
	// }

	// // Next, determine the value of the current_version_number custom option
	// // for the ProtoVersion message.
	// newVersionNumber := GetCurrentProtoVersion()

	// // Compare this to the currently deployed version number, which is written to a text file
	// // during the build-and-deploy Github Actions workflow.
	// deployedVersionNumber, err := readExpectedProtoVersion("../../.deployedprotoversion", false)

	// if err != nil {
	// 	t.Fatal("FAILURE LOADING .deployedprotoversion FILE!")
	// }

	// if breakingChangeDetected && (newVersionNumber == deployedVersionNumber) {
	// 	t.Fatalf("Breaking proto change detected without corresponding API version increase!\n%s\nEither fix your proto change so that it remains backwards-compatible, or increment the current_version_number option within the ProtoVersion message in proto/common.proto.\nIf you choose the latter route, then you must also modify the updateProtoVersion() method of any affected classes in the front-end code!\nThis ensures that saved browser data and old sim links will be automatically migrated to the new schema.", out)
	// } else if !breakingChangeDetected && (newVersionNumber != deployedVersionNumber) {
	// 	t.Fatal("API version increase detected without any breaking changes to protos!\nIf your proto changes are indeed backwards-compatible as detected, then revert the current_version_number option that you incremented in proto/common.proto back to its old value.\nIf your proto changes do in fact break saved browser data or old sim links, then make the breakage more explicit, such as by renaming an affected field.")
	// }

	// // If the above checks passed, then read in the actual results file for this test, which
	// // stores the local (pre-deployment) value of the version number. If necessary, force the
	// // developer to run make update-tests to keep this file in sync with proto/common.proto .
	// // This results file will be automatically copied into .deployedprotoversion as part of
	// // the build-and-deploy workflow.
	// localVersionNumber, err := readExpectedProtoVersion("TestProtoVersioning.results", false)

	// if err != nil {
	// 	t.Fatal("FAILURE LOADING TestProtoVersioning.results FILE!")
	// }

	// if localVersionNumber != newVersionNumber {
	// 	t.Logf("API version numbers in proto/common.proto and sim/core/TestProtoVersioning.results do not match: expected %d but was %d.\nIf you intentionally incremented the current_version_number option within the ProtoVersion message due to a breaking proto change, then the mismatch here is normal.\nIf this is the case, then simply update the results file by running make update-tests.\nThis will allow the test to pass, and also ensures that the .deployedprotoversion file is automatically updated during deployment.", localVersionNumber, newVersionNumber)
	// 	t.Fail()
	// }

	// // Write the new value to a .results.tmp file so that make update-tests
	// // will copy it automatically.
	// newVersionMessage := &proto.ProtoVersion{
	// 	SavedVersionNumber: newVersionNumber,
	// }
	// messageStr := prototext.Format(newVersionMessage)
	// messageStr = strings.ReplaceAll(messageStr, "  ", " ")
	// messageData := []byte(messageStr)
	// err = os.WriteFile("TestProtoVersioning.results.tmp", messageData, 0644)

	// if err != nil {
	// 	panic(err)
	// }
}
