package command

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"

	"github.com/gilcrest/diy-go-api/datastore"
	"github.com/gilcrest/diy-go-api/domain/errs"
	"github.com/gilcrest/diy-go-api/domain/logger"
	"github.com/gilcrest/diy-go-api/domain/secure"
	"github.com/gilcrest/diy-go-api/domain/secure/random"
	"github.com/gilcrest/diy-go-api/service"
)

// Genesis command runs the Genesis service and seeds the database.
func Genesis() (err error) {
	var (
		flgs        flags
		minlvl, lvl zerolog.Level
		ek          *[32]byte
	)

	// newFlags will retrieve the database info from the environment using ff
	flgs, err = newFlags([]string{"server"})
	if err != nil {
		return err
	}

	// determine minimum logging level based on flag input
	minlvl, err = zerolog.ParseLevel(flgs.logLvlMin)
	if err != nil {
		return err
	}

	// determine logging level based on flag input
	lvl, err = zerolog.ParseLevel(flgs.loglvl)
	if err != nil {
		return err
	}

	// setup logger with appropriate defaults
	lgr := logger.NewLogger(os.Stdout, minlvl, true)

	// logs will be written at the level set in NewLogger (which is
	// also the minimum level). If the logs are to be written at a
	// different level than the minimum, use SetGlobalLevel to set
	// the global logging level to that. Minimum rules will still
	// apply.
	if minlvl != lvl {
		zerolog.SetGlobalLevel(lvl)
	}

	// set global logging time field format to Unix timestamp
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lgr.Info().Msgf("minimum accepted logging level set to %s", minlvl)
	lgr.Info().Msgf("logging level set to %s", lvl)

	// set global to log errors with stack (or not) based on flag
	logger.WriteErrorStackGlobal(flgs.logErrorStack)
	lgr.Info().Msgf("log error stack global set to %t", flgs.logErrorStack)

	if flgs.encryptkey == "" {
		lgr.Fatal().Msg("no encryption key found")
	}

	// decode and retrieve encryption key
	ek, err = secure.ParseEncryptionKey(flgs.encryptkey)
	if err != nil {
		lgr.Fatal().Err(err).Msg("secure.ParseEncryptionKey() error")
	}

	ctx := context.Background()

	// initialize PostgreSQL database
	var (
		dbpool  *pgxpool.Pool
		cleanup func()
	)
	dbpool, cleanup, err = datastore.NewPostgreSQLPool(ctx, newPostgreSQLDSN(flgs), lgr)
	if err != nil {
		lgr.Fatal().Err(err).Msg("datastore.NewPostgreSQLPool error")
	}
	defer cleanup()

	s := service.GenesisService{
		Datastorer:            datastore.NewDatastore(dbpool),
		RandomStringGenerator: random.CryptoGenerator{},
		EncryptionKey:         ek,
	}

	var b []byte
	b, err = os.ReadFile(genesisRequestFile)
	if err != nil {
		return errs.E(err)
	}
	f := service.GenesisRequest{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return errs.E(err)
	}

	var response service.GenesisResponse
	response, err = s.Seed(ctx, &f)
	if err != nil {
		return err
	}

	var responseJSON []byte
	responseJSON, err = json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(service.LocalJSONGenesisResponseFile, responseJSON, 0644)
	if err != nil {
		return err
	}

	fmt.Println(string(responseJSON))

	return nil
}

// NewEncryptionKey generates a random 256-bit key and prints it to standard out.
// It will return an error if the system's secure random number generator fails
// to function correctly, in which case the caller should not continue.
// Taken from https://github.com/gtank/cryptopasta/blob/master/encrypt.go
func NewEncryptionKey() {
	lgr := logger.NewLogger(os.Stdout, zerolog.DebugLevel, true)

	keyBytes, err := secure.NewEncryptionKey()
	if err != nil {
		lgr.Fatal().Err(err).Msg("EncryptionKey() error")
	}

	fmt.Printf("Key Ciphertext:\t[%s]\n", hex.EncodeToString(keyBytes[:]))
}
