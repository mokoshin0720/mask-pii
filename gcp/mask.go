package gcp

import (
	"context"
	"fmt"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"github.com/mokoshin0720/mask-pii/gcp/config"
	"google.golang.org/api/option"
)

// Mask inputで指定された文字列にinfoTypeNamesで指定された情報タイプのマスク処理を行う
func Mask(input string) (string, error) {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx, option.WithAPIKey(config.Config.GcpApiKey))
	if err != nil {
		return "", fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		InspectConfig: &dlppb.InspectConfig{}, // NOTE: 全ての情報タイプを検出対象とする
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{}, // Match all info types.
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CharacterMaskConfig{
									CharacterMaskConfig: &dlppb.CharacterMaskConfig{},
								},
							},
						},
					},
				},
			},
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
	}

	r, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return "", fmt.Errorf("DeidentifyContent: %w", err)
	}

	return r.GetItem().GetValue(), nil
}
