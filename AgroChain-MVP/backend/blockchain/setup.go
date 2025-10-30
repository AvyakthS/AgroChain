package blockchain

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type FabricService struct {
	Gateway  *client.Gateway
	Contract *client.Contract
}

func Initialize() (*FabricService, error) {
	peerEndpoint := os.Getenv("PEER_ENDPOINT")
	gatewayPeer := os.Getenv("GATEWAY_PEER")
	mspID := os.Getenv("MSP_ID")
	tlsCertPath := os.Getenv("TLS_CERT_PATH")
	certPath := os.Getenv("CERT_PATH")
	keyPathDir := os.Getenv("KEY_PATH_DIR")
	channelName := os.Getenv("CHANNEL_NAME")
	chaincodeName := os.Getenv("CHAINCODE_NAME")

	clientConn, err := newGrpcConnection(tlsCertPath, peerEndpoint, gatewayPeer)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	id, err := newIdentity(certPath, mspID)
	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	sign, err := newSign(keyPathDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create sign: %w", err)
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	fmt.Println("Successfully connected to the blockchain gateway.")
	return &FabricService{Gateway: gw, Contract: contract}, nil
}

func (s *FabricService) Close() {
	if s.Gateway != nil {
		s.Gateway.Close()
	}
}

// --- THIS FUNCTION IS NOW CORRECTED ---
func newGrpcConnection(tlsCertPath, peerEndpoint, gatewayPeer string) (*grpc.ClientConn, error) {
	tlsCertPEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLS cert file: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(tlsCertPEM) {
		return nil, fmt.Errorf("failed to add TLS cert to pool")
	}

	config := &tls.Config{
		RootCAs:    certPool,
		ServerName: gatewayPeer,
	}
	transportCredentials := credentials.NewTLS(config)

	return grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
}
// --- END OF CORRECTION ---

func newIdentity(certPath, mspID string) (*identity.X509Identity, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user cert file: %w", err)
	}
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}
	return identity.NewX509Identity(mspID, cert)
}

func newSign(keyPathDir string) (identity.Sign, error) {
	files, err := os.ReadDir(keyPathDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore directory: %w", err)
	}
	privateKeyPath := path.Join(keyPathDir, files[0].Name())

	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	return identity.NewPrivateKeySign(privateKey)
}