build-sam:
	sam build

debug-with-sam: build-sam
	sam local invoke


