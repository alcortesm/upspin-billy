package upspinfs

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v2/test"
	"upspin.io/bind"
	"upspin.io/client"
	"upspin.io/config"
	dirserver "upspin.io/dir/inprocess"
	"upspin.io/errors"
	"upspin.io/factotum"
	keyserver "upspin.io/key/inprocess"
	storeserver "upspin.io/store/inprocess"
	"upspin.io/test/testutil"
	"upspin.io/upspin"
)

func Test(t *testing.T) { TestingT(t) }

type UpspinSuite struct {
	test.FilesystemSuite
	cfg upspin.Config
}

var _ = Suite(&UpspinSuite{})

// fixtures
var (
	userName = upspin.UserName("user1@google.com")
	// from $GOPATH/src/upspin.io/key/testdata/user1/
	publicKey = upspin.PublicKey(`p256
104278369061367353805983276707664349405797936579880352274235000127123465616334
26941412685198548642075210264642864401950753555952207894712845271039438170192
`)
	// We will use in-memory, non-persistent key, dir and store servers.
	// The InProcess transport is just perfect for them; it ignores addresses,
	// this is, all servers will be at the same endpoint:
	serverEP = upspin.Endpoint{Transport: upspin.InProcess}
)

func (s *UpspinSuite) SetUpSuite(c *C) {
	s.setUpClientConfig(c)
	s.assertConfig(c)

	setUpServers(c, s.cfg)
	s.assertUser(c)
	s.assertEmptyHome(c)

	s.FilesystemSuite.FS = New(client.New(s.cfg), userName)
}

// Initialize the client config in its receiver to the fixture values.
func (s *UpspinSuite) setUpClientConfig(c *C) {
	s.cfg = config.New()
	s.cfg = config.SetPacking(s.cfg, upspin.EEIntegrityPack)
	s.cfg = config.SetUserName(s.cfg, userName)
	s.cfg = config.SetKeyEndpoint(s.cfg, serverEP)
	s.cfg = config.SetStoreEndpoint(s.cfg, serverEP)
	s.cfg = config.SetDirEndpoint(s.cfg, serverEP)

	f, err := factotum.NewFromDir(testutil.Repo("key", "testdata", "user1"))
	c.Assert(err, IsNil)
	s.cfg = config.SetFactotum(s.cfg, f)
}

// Runs and registers key, dir, and store servers for a the user
// described in cfg.  Its home is initialized empty.
func setUpServers(c *C, cfg upspin.Config) {
	err := bind.RegisterKeyServer(upspin.InProcess, keyserver.New())
	c.Assert(err, IsNil)
	err = bind.RegisterStoreServer(upspin.InProcess, storeserver.New())
	c.Assert(err, IsNil)
	err = bind.RegisterDirServer(upspin.InProcess, dirserver.New(cfg))
	c.Assert(err, IsNil)

	user := &upspin.User{
		Name:      cfg.UserName(),
		Dirs:      []upspin.Endpoint{cfg.DirEndpoint()},
		Stores:    []upspin.Endpoint{cfg.StoreEndpoint()},
		PublicKey: cfg.Factotum().PublicKey(),
	}

	// add user to key server
	key, err := bind.KeyServer(cfg, cfg.KeyEndpoint())
	c.Assert(err, IsNil)
	err = key.Put(user)
	c.Assert(err, IsNil)

	homePath := upspin.PathName(cfg.UserName() + "/")
	entry := &upspin.DirEntry{
		Name:       homePath,
		SignedName: homePath,
		Attr:       upspin.AttrDirectory,
		Writer:     cfg.UserName(),
	}

	// add a home directory entry for the user at the dir server
	dir, err := bind.DirServer(cfg, cfg.DirEndpoint())
	c.Assert(err, IsNil)
	_, err = dir.Put(entry)
	c.Assert(err, IsNil)
}

func (s *UpspinSuite) assertConfig(c *C) {
	c.Assert(s.cfg.UserName(), Equals, userName)
	c.Assert(s.cfg.Factotum().PublicKey(), Equals, publicKey)
	c.Assert(s.cfg.KeyEndpoint(), Equals, serverEP)
	c.Assert(s.cfg.DirEndpoint(), Equals, serverEP)
	c.Assert(s.cfg.StoreEndpoint(), Equals, serverEP)
}

func (s *UpspinSuite) assertUser(c *C) {
	e := s.cfg.KeyEndpoint()
	c.Assert(e.Transport, Equals, upspin.InProcess)
	keyServer, err := bind.KeyServer(s.cfg, e)
	c.Assert(err, IsNil)

	user, err := keyServer.Lookup(userName)
	c.Assert(err, IsNil)
	c.Assert(user.Name, Equals, userName)
	c.Assert(user.Dirs, DeepEquals, []upspin.Endpoint{serverEP})
	c.Assert(user.Stores, DeepEquals, []upspin.Endpoint{serverEP})
	c.Assert(user.PublicKey, Equals, publicKey)
}

func (s *UpspinSuite) assertEmptyHome(c *C) {
	e := s.cfg.DirEndpoint()
	c.Assert(e.Transport, Equals, upspin.InProcess)
	dirServer, err := bind.DirServer(s.cfg, e)
	c.Assert(err, IsNil)

	path := upspin.PathName(userName + "/")
	de, err := dirServer.Lookup(path)
	c.Assert(err, IsNil)
	c.Assert(de.IsDir(), Equals, true)
	c.Assert(de.Writer, Equals, userName)
	c.Assert(de.Name, Equals, path)
	c.Assert(len(de.Blocks), Equals, 1)
	block := de.Blocks[0]
	c.Assert(block.Offset, Equals, int64(0))
	c.Assert(block.Size, Equals, int64(0))
	c.Assert(block.Location.Endpoint, Equals, s.cfg.StoreEndpoint())

	e = s.cfg.StoreEndpoint()
	c.Assert(e.Transport, Equals, upspin.InProcess)
	storeServer, err := bind.StoreServer(s.cfg, e)
	c.Assert(err, IsNil)

	data, refData, locations, err := storeServer.Get(block.Location.Reference)
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte{})
	c.Assert(refData.Reference, Equals, block.Location.Reference)
	c.Assert(refData.Volatile, Equals, false)
	c.Assert(refData.Duration, Equals, time.Duration(0))
	c.Assert(locations, IsNil)
}

func (s *UpspinSuite) SetUpTest(c *C) {
	s.assertEmptyHome(c)
}

func (s *UpspinSuite) TearDownTest(c *C) {
	dirServer, err := bind.DirServerFor(s.cfg, "")
	c.Assert(err, IsNil)
	deleteAll(dirServer, upspin.PathName(s.cfg.UserName()+"/"))
}

var (
	errNotExist = errors.E(errors.NotExist)
)

// deleteAll recursively deletes the directory named by path through the
// provided DirServer, first deleting path/Access and then path/*.
func deleteAll(dir upspin.DirServer, path upspin.PathName) error {
	if _, err := dir.Delete(path + "/Access"); err != nil {
		if !errors.Match(errNotExist, err) {
			return err
		}
	}
	entries, err := dir.Glob(string(path + "/*"))
	if err != nil && err != upspin.ErrFollowLink {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			if err := deleteAll(dir, e.Name); err != nil {
				return err
			}
		}
		if _, err := dir.Delete(e.Name); err != nil {
			return err
		}
	}
	return nil
}
