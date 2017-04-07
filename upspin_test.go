package upspinfs

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
	"upspin.io/bind"
	"upspin.io/config"
	dirserver "upspin.io/dir/inprocess"
	"upspin.io/factotum"
	keyserver "upspin.io/key/inprocess"
	storeserver "upspin.io/store/inprocess"
	"upspin.io/test/testutil"
	"upspin.io/upspin"
)

func Test(t *testing.T) { TestingT(t) }

type UpspinSuite struct {
	//test.FilesystemSuite
	userName upspin.UserName
	cfg      upspin.Config
}

var _ = Suite(&UpspinSuite{})

func (s *UpspinSuite) SetUpSuite(c *C) {
	s.userName = upspin.UserName("user1@google.com")
	s.cfg = inprocessConfig(s.userName, "")
}

func (s *UpspinSuite) SetUpTest(c *C) {
	//client := client.New(s.cfg)
	//s.FilesystemSuite.Fs = New(client, s.userName)
}

func inprocessConfig(userName upspin.UserName, publicKey upspin.PublicKey) upspin.Config {
	cfg := baseConfig()
	cfg = config.SetUserName(cfg, userName)
	key, _ := bind.KeyServer(cfg, cfg.KeyEndpoint())
	checkTransport(key)
	dir, _ := bind.DirServer(cfg, cfg.DirEndpoint())
	checkTransport(dir)
	if cfg.Factotum().PublicKey() == "" {
		panic("empty public key")
	}
	user := &upspin.User{
		Name:      userName,
		Dirs:      []upspin.Endpoint{cfg.DirEndpoint()},
		Stores:    []upspin.Endpoint{cfg.StoreEndpoint()},
		PublicKey: cfg.Factotum().PublicKey(),
	}
	err := key.Put(user)
	if err != nil {
		panic(err)
	}
	name := upspin.PathName(userName) + "/"
	entry := &upspin.DirEntry{
		Name:       name,
		SignedName: name,
		Attr:       upspin.AttrDirectory,
		Writer:     userName,
	}
	_, err = dir.Put(entry)
	if err != nil {
		panic(err)
	}
	return cfg
}

func checkTransport(s upspin.Service) {
	if s == nil {
		panic(fmt.Sprintf("nil service"))
	}
	if t := s.Endpoint().Transport; t != upspin.InProcess {
		panic(fmt.Sprintf("bad transport %v, want inprocess", t))
	}
}

func baseConfig() upspin.Config {
	inProcess := upspin.Endpoint{
		Transport: upspin.InProcess,
		NetAddr:   "", // ignored
	}

	f, err := factotum.NewFromDir(testutil.Repo("key", "testdata", "user1")) // Always use user1's keys.
	if err != nil {
		panic("cannot initialize factotum: " + err.Error())
	}

	cfg := config.New()
	cfg = config.SetPacking(cfg, upspin.EEIntegrityPack)
	cfg = config.SetKeyEndpoint(cfg, inProcess)
	cfg = config.SetStoreEndpoint(cfg, inProcess)
	cfg = config.SetDirEndpoint(cfg, inProcess)
	cfg = config.SetFactotum(cfg, f)

	bind.RegisterKeyServer(upspin.InProcess, keyserver.New())
	bind.RegisterStoreServer(upspin.InProcess, storeserver.New())
	bind.RegisterDirServer(upspin.InProcess, dirserver.New(cfg))

	return cfg
}

func (s *UpspinSuite) TearDownTest(c *C) {
}

func (s *UpspinSuite) TestSetUpSuite(c *C) {
	keyServer, err := bind.KeyServer(s.cfg, s.cfg.KeyEndpoint())
	c.Assert(err, IsNil)

	userName := upspin.UserName("user1@google.com")
	user, err := keyServer.Lookup(userName)
	c.Assert(err, IsNil)

	c.Assert(user.Name, Equals, userName)
	fmt.Println(user.Dirs)
	fmt.Println(user.Stores)
	fmt.Println(user.PublicKey)

}
