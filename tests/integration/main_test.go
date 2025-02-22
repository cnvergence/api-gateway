package api_gateway

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func InitializeIstioJwtTests(ctx *godog.TestSuiteContext) {
	InitializeScenarioIstioJWT(ctx.ScenarioContext())
}

func TestIstioJwt(t *testing.T) {
	InitTestSuite()

	orgJwtHandler, err := SwitchJwtHandler("istio")
	if err != nil {
		log.Print(err.Error())
		t.Fatalf("unable to switch to Istio jwtHandler")
	}

	SetupCommonResources("istio-jwt")

	apiGatewayIstioJwtOpts := goDogOpts
	apiGatewayIstioJwtOpts.Paths = []string{"features/istio-jwt/istio_jwt.feature"}
	apiGatewayIstioJwtOpts.Concurrency = conf.TestConcurency

	apiGatewayIstioJwtSuite := godog.TestSuite{
		Name:                 "istio-jwt",
		TestSuiteInitializer: InitializeIstioJwtTests,
		Options:              &apiGatewayIstioJwtOpts,
	}

	defer cleanUp(orgJwtHandler)

	testExitCode := apiGatewayIstioJwtSuite.Run()
	if testExitCode != 0 {
		t.Fatalf("non-zero status returned, failed to run feature tests, Pod list: %s\n APIRules: %s\n", getPodListReport(), getApiRules())
	}
}

func cleanUp(orgJwtHandler string) {
	res := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	err := k8sClient.Resource(res).Delete(context.Background(), namespace, v1.DeleteOptions{})

	if err != nil {
		log.Print(err.Error())
	}

	if os.Getenv(exportResultVar) == "true" {
		generateReport()
	}

	_, err = SwitchJwtHandler(orgJwtHandler)
	if err != nil {
		log.Print(err.Error())
		t.Fatalf("unable to switch back to original jwtHandler")
	}
}
