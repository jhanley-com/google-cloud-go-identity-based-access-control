@REM This code gets the Project ID from gcloud
call gcloud config get-value project > project.tmp
for /f "delims=" %%x in (project.tmp) do set GCP_PROJECT_ID=%%x
echo Project ID: %GCP_PROJECT_ID%

del project.tmp

@echo on

set GCP_SA_1=first-service-account@%GCP_PROJECT_ID%.iam.gserviceaccount.com
set GCP_SA_2=second-service-account@%GCP_PROJECT_ID%.iam.gserviceaccount.com

set GCP_SA_FILE_1=first-service-account.json
set GCP_SA_FILE_2=second-service-account.json

set GCS_BUCKET_NAME=%GCP_PROJECT_ID%-xtest

set GCS_BUCKET_ROLE=legacyBucketReader
set GCS_OBJECT_ROLE=legacyObjectReader

call gcloud iam service-accounts create first-service-account ^
--display-name="First Service Account"
@echo on

call gcloud iam service-accounts keys create %GCP_SA_FILE_1% ^
--iam-account="%GCP_SA_1%" ^
--key-file-type=json
@echo on

call gcloud iam service-accounts create second-service-account ^
--display-name="Second Service Account"
@echo on

call gcloud iam service-accounts keys create %GCP_SA_FILE_2% ^
--iam-account="%GCP_SA_2%" ^
--key-file-type=json
@echo on

call gcloud projects add-iam-policy-binding %GCP_PROJECT_ID% ^
--member serviceAccount:"%GCP_SA_2%" ^
--role roles/storage.objectViewer
@echo on

call gsutil mb gs://%GCS_BUCKET_NAME%
@echo on

call gsutil cp %GCP_SA_FILE_2% gs://%GCS_BUCKET_NAME%
@echo on

@REM gsutil iam ch serviceAccount:%GCP_SA_1%:%GCS_BUCKET_ROLE% gs://%GCS_BUCKET_NAME%/
@echo on

gsutil iam ch serviceAccount:%GCP_SA_1%:%GCS_OBJECT_ROLE% gs://%GCS_BUCKET_NAME%/%GCP_SA_FILE_2%
@echo on
