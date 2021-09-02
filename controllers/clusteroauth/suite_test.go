// Copyright Red Hat

package clusteroauth

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	idpclientset "github.com/identitatem/idp-client-api/api/client/clientset/versioned"
	identitatemv1alpha1 "github.com/identitatem/idp-client-api/api/identitatem/v1alpha1"
	idpconfig "github.com/identitatem/idp-client-api/config"
	openshiftconfigv1 "github.com/openshift/api/config/v1"
	clientsetcluster "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clientsetwork "open-cluster-management.io/api/client/work/clientset/versioned"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	clusterv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
	workv1 "open-cluster-management.io/api/work/v1"
	clusteradmasset "open-cluster-management.io/clusteradm/pkg/helpers/asset"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var clientSetMgmt *idpclientset.Clientset
var clientSetStrategy *idpclientset.Clientset
var clientSetCluster *clientsetcluster.Clientset
var clientSetWork *clientsetwork.Clientset
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"ClusterOAuth Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	readerIDP := idpconfig.GetScenarioResourcesReader()
	clusterOAuthCRD, err := getCRD(readerIDP, "crd/bases/identityconfig.identitatem.io_clusteroauths.yaml")
	Expect(err).Should(BeNil())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDs: []client.Object{
			clusterOAuthCRD,
		},
		CRDDirectoryPaths: []string{
			//DV added this line and copyed the authrealms CRD
			filepath.Join("..", "..", "test", "config", "crd", "external"),
		},
		ErrorIfCRDPathMissing: true,
	}
	err = identitatemv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = corev1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = clusterv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = clusterv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = workv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = openshiftconfigv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	clientSetMgmt, err = idpclientset.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())
	Expect(clientSetMgmt).ToNot(BeNil())

	clientSetStrategy, err = idpclientset.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())
	Expect(clientSetStrategy).ToNot(BeNil())

	clientSetCluster, err = clientsetcluster.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())
	Expect(clientSetCluster).ToNot(BeNil())

	clientSetWork, err = clientsetwork.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())
	Expect(clientSetWork).ToNot(BeNil())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	By("Creating infra", func() {
		infraConfig := &openshiftconfigv1.Infrastructure{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cluster",
			},
			Spec: openshiftconfigv1.InfrastructureSpec{},
			Status: openshiftconfigv1.InfrastructureStatus{
				APIServerURL: "http://127.0.0.1:6443",
			},
		}
		err := k8sClient.Create(context.TODO(), infraConfig)
		Expect(err).NotTo(HaveOccurred())
	})

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Process clusterOAuth for Strategy backplane: ", func() {
	AuthRealmName := "my-authrealm"
	//AuthRealmNameSpace := "my-authrealmns"
	//CertificatesSecretRef := "my-certs"
	StrategyName := AuthRealmName + "-backplane"
	PlacementStrategyName := StrategyName
	// PlacementName := AuthRealmName
	ClusterName := "my-cluster"
	MyIDPName := "my-idp"

	It("process a ClusterOAuth CR", func() {
		By(fmt.Sprintf("creation of cluster namespace %s", ClusterName), func() {
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: ClusterName,
				},
			}
			err := k8sClient.Create(context.TODO(), ns)
			Expect(err).To(BeNil())
		})

		var secret *corev1.Secret
		By(fmt.Sprintf("creation of IDP secret in cluster namespace %s", ClusterName), func() {
			secret = &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: corev1.SchemeGroupVersion.String(),
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MyIDPName,
					Namespace: ClusterName,
				},
			}
			err := k8sClient.Create(context.TODO(), secret)
			Expect(err).To(BeNil())

			//			Expect(secret.TypeMeta.Kind).To(Not(BeEmpty()))
			//			Expect(secret.TypeMeta.APIVersion).To(Not(BeEmpty()))

			secretAfter := &corev1.Secret{}
			err = k8sClient.Get(context.TODO(), types.NamespacedName{Name: MyIDPName, Namespace: ClusterName}, secretAfter)
			Expect(err).To(BeNil())

		})

		By(fmt.Sprintf("creation of ClusterOAuth for mangaed cluster %s", ClusterName), func() {
			clusterOAuth := &identitatemv1alpha1.ClusterOAuth{
				TypeMeta: metav1.TypeMeta{
					APIVersion: identitatemv1alpha1.SchemeGroupVersion.String(),
					Kind:       "ClusterOAuth",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      PlacementStrategyName,
					Namespace: ClusterName,
				},
				Spec: identitatemv1alpha1.ClusterOAuthSpec{
					OAuth: &openshiftconfigv1.OAuth{
						TypeMeta: metav1.TypeMeta{
							APIVersion: openshiftconfigv1.SchemeGroupVersion.String(),
							Kind:       "OAuth",
						},

						Spec: openshiftconfigv1.OAuthSpec{
							IdentityProviders: []openshiftconfigv1.IdentityProvider{
								{
									Name:          MyIDPName,
									MappingMethod: openshiftconfigv1.MappingMethodClaim,
									IdentityProviderConfig: openshiftconfigv1.IdentityProviderConfig{
										Type: openshiftconfigv1.IdentityProviderTypeGitHub,
										GitHub: &openshiftconfigv1.GitHubIdentityProvider{
											ClientID: "me",
											ClientSecret: openshiftconfigv1.SecretNameReference{
												Name: secret.Name,
											},
										},
									},
								},
							},
						},
					},
				},
			}
			err := k8sClient.Create(context.TODO(), clusterOAuth)
			Expect(err).To(BeNil())
		})

		By("Calling reconcile", func() {
			r := &ClusterOAuthReconciler{
				Client: k8sClient,
				Log:    logf.Log,
				Scheme: scheme.Scheme,
			}
			req := ctrl.Request{}
			req.Name = PlacementStrategyName
			req.Namespace = ClusterName
			_, err := r.Reconcile(context.TODO(), req)
			Expect(err).To(BeNil())
		})
		By("Checking manifestwork", func() {
			mw := &workv1.ManifestWork{}
			err := k8sClient.Get(context.TODO(), types.NamespacedName{Name: "idp-backplane", Namespace: ClusterName}, mw)
			//var mw *workv1.ManifestWork
			//mw, err := clientSetWork.WorkV1().ManifestWorks(ClusterName).Get(context.TODO(), "idp-backplane", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(mw.Name).To(Equal("idp-backplane"))
			Expect(mw.Namespace).To(Equal(ClusterName))
			// should find manifest for OAuth and manifest for Secret
			Expect(len(mw.Spec.Workload.Manifests)).To(Equal(2))
		})
	})
})

func getCRD(reader *clusteradmasset.ScenarioResourcesReader, file string) (*apiextensionsv1.CustomResourceDefinition, error) {
	b, err := reader.Asset(file)
	if err != nil {
		return nil, err
	}
	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := yaml.Unmarshal(b, crd); err != nil {
		return nil, err
	}
	return crd, nil
}
