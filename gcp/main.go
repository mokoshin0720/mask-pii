import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// mask deidentifies the input by masking all provided info types with maskingCharacter
// and prints the result to w.
func mask(w io.Writer, projectID, input string, infoTypeNames []string, maskingCharacter string, numberToMask int32) error {
	// projectID := "my-project-id"
	// input := "My SSN is 111222333"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}
	// maskingCharacter := "+"
	// numberToMask := 6
	// Will print "My SSN is ++++++333"

	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()
	// Convert the info type strings to a list of InfoTypes.
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}
	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: infoTypes,
		},
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{}, // Match all info types.
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CharacterMaskConfig{
									CharacterMaskConfig: &dlppb.CharacterMaskConfig{
										MaskingCharacter: maskingCharacter,
										NumberToMask:     numberToMask,
									},
								},
							},
						},
					},
				},
			},
		},
		// The item to analyze.
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
	}
	// Send the request.
	r, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return fmt.Errorf("DeidentifyContent: %w", err)
	}
	// Print the result.
	fmt.Fprint(w, r.GetItem().GetValue())
	return nil
}
