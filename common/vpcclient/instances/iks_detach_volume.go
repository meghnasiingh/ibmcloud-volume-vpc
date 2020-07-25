/**
 * Copyright 2020 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instances

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/client"
	"github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/vpc/vpcclient/models"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// DetachVolume retrives the volume attach status with givne volume attachment details
func (vs *IKSVolumeAttachService) DetachVolume(volumeAttachmentTemplate *models.VolumeAttachment, ctxLogger *zap.Logger) (*http.Response, error) {
	defer util.TimeTracker("IKS DetachVolume", time.Now())

	operation := &client.Operation{
		Name:        "DetachVolume",
		Method:      "DELETE",
		PathPattern: vs.pathPrefix + "deleteAttachment",
	}

	apiErr := vs.receiverError

	operationRequest := vs.client.NewRequest(operation)
	operationRequest = operationRequest.SetQueryValue(IksClusterQueryKey, *volumeAttachmentTemplate.ClusterID)
	operationRequest = operationRequest.SetQueryValue(IksWorkerQueryKey, *volumeAttachmentTemplate.InstanceID)
	operationRequest = operationRequest.SetQueryValue(IksVolumeAttachmentIDQueryKey, volumeAttachmentTemplate.ID)

	ctxLogger.Info("Equivalent curl command and query parameters", zap.Reflect("URL", operationRequest.URL()), zap.Reflect("volumeAttachmentTemplate", volumeAttachmentTemplate), zap.Reflect("Operation", operation), zap.Reflect(IksClusterQueryKey, *volumeAttachmentTemplate.ClusterID), zap.Reflect(IksWorkerQueryKey, *volumeAttachmentTemplate.InstanceID), zap.Reflect(IksVolumeAttachmentIDQueryKey, volumeAttachmentTemplate.ID))

	resp, err := operationRequest.JSONError(apiErr).Invoke()
	if err != nil {
		ctxLogger.Error("Error occured while deleting volume attachment", zap.Error(err))
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			// volume attachment is deleted, no need to retry
			return resp, apiErr
		}
	}

	ctxLogger.Info("Successfuly deleted the volume attachment")
	return resp, err
}