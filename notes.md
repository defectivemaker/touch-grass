AES-256 GCM is used to encrypt the data. The client and server has the same key. Since this project is now open source, I ommited the key. Variable KeyOne in client/common/common.go and server/internal/common/common.go

I have removed all certificates that are stored here, you will have to generate them and import them in yourself. The server certs will go in server/internal/server/local-certs

I have replaced any mention of touchgrass.au with YOUR_DOMAIN. Grep this out and replace it with your own domain if you plan on making this yourself
