package handlers

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"ussd_ethereum/internal/eth"
	"ussd_ethereum/internal/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofiber/fiber/v2"

	"github.com/AndroidStudyOpenSource/africastalking-go/sms"
)

const (
	apiKey    = "atsk_35079f792da37a72be09065f7659c8de5973dd655160c7280e7268e12dcbb0f0760fa898"
	username  = "sandbox"
	shortCode = "44472"
	env       = ""
)

var menuTree *MenuTree

var smsService = sms.NewService(apiKey, username, env)

type Data struct {
	PhoneNumber string
	Address     string
}

func (h *Handler) CallbackHandler(c *fiber.Ctx) error {

	text := c.FormValue("text")
	phoneNumber := c.FormValue("phoneNumber")

	steps := strings.Split(text, "*")

	if len(steps) > 1 && steps[len(steps)-1] == "0" {
		steps = steps[:len(steps)-2]
	}

	if len(steps) > 1 && steps[len(steps)-1] == "0" {
		steps = steps[:len(steps)-2]
	}

	currentText := strings.Join(steps, "*")

	var response string
	fmt.Println(steps)

	switch len(steps) {
	case 1:
		if currentText == "" {
			// Main menu
			response = "CON Welcome to the ETH Service\n1. Create Account\n2. Account Details\n3. Send ETH\n4. Receive ETH\n"
		} else {
			ksRecord, err := h.DB.SelectWalletByPhone(phoneNumber)
			if err != nil {
				fmt.Println(err)
			}

			switch steps[0] {
			case "1":
				if ksRecord != nil {
					response = fmt.Sprintf("END Account for %s already exists.", phoneNumber)
				} else {
					response = "CON Set PIN:\n0. Back"
				}

			case "2":
				if ksRecord == nil || err == sql.ErrNoRows {
					response = fmt.Sprintf("END Account for %s does not exist.", phoneNumber)
				} else {
					response = "CON Enter your PIN to view account details:\n0. Back"
				}
			case "3":
				response = "CON Enter recipient's phone number:\n0. Back"
			case "4":
				if ksRecord == nil || err == sql.ErrNoRows {
					response = fmt.Sprintf("END Account for %s does not exist\n", phoneNumber)
					response += "Create Account"
				} else {
					response = fmt.Sprintf("END Phone number: %s\nPublic key: \n%s", phoneNumber, ksRecord.PublicKey)
				}
			default:
				response = "END Invalid choice. Try again."
			}
		}
	case 2:
		switch steps[0] {
		case "1": // Create Account
			pin := steps[1]

			if len(pin) < 4 {
				response = "CON Invalid PIN length:\n0. Back"
			} else {
				ks, err := eth.NewKeyStore(pin)
				if err != nil {
					response = fmt.Sprintf("END Account creation failed %s.", err)
					return err
				}
				addr := ks.Address.Hex()

				err = h.DB.InsertWallet(phoneNumber, pin, addr, ks.URL.Path)
				if err != nil {
					fmt.Println(err)
				}
				response = fmt.Sprintf("END Account created successfully with address %s.", addr)
			}
		case "2": // Account Details - PIN entry
			pin := steps[1]
			ksRecord, err := h.DB.SelectWalletByPhone(phoneNumber)
			if err != nil {
				response = fmt.Sprintf("END Account details failed %s.", err)
				return err
			}
			validPass := ksRecord.Pin
			if pin == validPass {
				response = "CON Select an option:\n1. View\n0. Back"
			} else {
				response = "END Invalid PIN. Access denied."
			}
		case "3": // Send ETH
			response = "CON Enter amount in KES:\n0. Back"
		default:
			response = "END Invalid choice."
		}
	case 3:
		ksRecord, err := h.DB.SelectWalletByPhone(phoneNumber)
		if err != nil {
			response = fmt.Sprintf("END Account details failed %s.", err)
			return err
		}
		switch steps[0] {
		case "2": // Account Details - After PIN authentication
			fmt.Println(ksRecord)
			address := common.HexToAddress(ksRecord.PublicKey)
			bal, err := h.Client.BalanceAt(c.Context(), address, nil)
			if err != nil {
				panic(err)
			}
			response = fmt.Sprintf("END Phone number: %s\nPublic Key: %s\nBalance: %f(ETH)", phoneNumber, ksRecord.PublicKey, utils.WeiToEth(bal))
		case "3": // Send ETH
			recipient := steps[1]
			amount := steps[2]
			recipientPhone, err := utils.FormatPhoneNumber(recipient)
			if err != nil {
				response = fmt.Sprintf("error: %s\n", err.Error())
				return err
			}

			recipientKS, err := h.DB.SelectWalletByPhone(recipientPhone)
			if recipientKS == nil || err == sql.ErrNoRows {
				response = fmt.Sprintf("Account for %s does not exist", recipient)
			} else {
				if amount == "" || strings.ContainsAny(amount, "abcdefghijklmnopqrstuvwxyz") {
					response = "CON Invalid amount. Please enter a valid amount in ETH:\n0. Back"
				} else {
					response = fmt.Sprintf("CON Enter your PIN to confirm sending %s KES to %s:\n0. Back", amount, recipient)
				}
			}

		case "5": // Buy Goods/Services
			till := steps[1]
			amount := steps[2]
			if amount == "" || strings.ContainsAny(amount, "abcdefghijklmnopqrstuvwxyz") {
				response = "CON Invalid amount. Please enter a valid amount in ETH:\n0. Back"
			} else {
				response = fmt.Sprintf("END You have paid %s ETH to Till Number %s.", amount, till)
			}
		default:
			response = "END Invalid choice."
		}
	case 4:
		switch steps[0] {
		case "2":
			response = "END Here"
		case "3": // Send ETH - PIN entry

			recipient := steps[1]
			amount := steps[2]
			pin := steps[3]

			ksRecord, err := h.DB.SelectWalletByPhone(phoneNumber)
			if err != nil {
				response = fmt.Sprintf("END Account details failed %s.", err)
				return err
			}
			validPass := ksRecord.Pin

			if pin == validPass {
				jsonBytes, err := os.ReadFile(ksRecord.KeystorePath)

				if err != nil {
					panic(err)
				}
				key, err := keystore.DecryptKey(jsonBytes, pin)
				privateKey := key.PrivateKey
				publicKey := privateKey.Public()
				publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
				if !ok {
					panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
				}

				fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
				nonce, err := h.Client.PendingNonceAt(context.Background(), fromAddress)
				if err != nil {
					panic(err)
				}

				amountInUSD, err := utils.ConvertKESToUSD(amount)
				amountInEth, err := utils.ConvertUSDToEth(amountInUSD)
				if err != nil {
					slog.Error("conversion error", "err", err)
				}
				slog.Info("amount in eth", "amount", amountInEth)

				// value := big.NewInt(1000000000000000000) // in wei (1 eth)
				weiAmt := amountInEth * 1000000000000000000
				value := big.NewInt(int64(weiAmt)) // in wei (1 eth)
				gasLimit := uint64(21000)          // in units
				gasPrice, err := h.Client.SuggestGasPrice(context.Background())
				if err != nil {
					panic(err)
				}
				recipientPhone, err := utils.FormatPhoneNumber(recipient)
				if err != nil {
					response = fmt.Sprintf("error: %s\n", err.Error())
					return err
				}
				recipientKS, err := h.DB.SelectWalletByPhone(recipientPhone)

				toAddress := common.HexToAddress(recipientKS.PublicKey)
				var data []byte
				tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

				chainID, err := h.Client.NetworkID(context.Background())
				if err != nil {
					panic(err)
				}

				signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
				if err != nil {
					panic(err)
				}

				err = h.Client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					panic(err)
				}

				fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

				response = fmt.Sprintf("END You have successfully sent %s ETH to %s.", amount, recipient)
				smsResponse, err := smsService.Send(phoneNumber, "+254716266543", "You have recieved")
				if err != nil {
					log.Panicf("failed to send sms: %v", err)
				}
				slog.Info("sms sent", "response", smsResponse)

			} else {
				response = "END Invalid PIN. Transaction failed."
			}
		default:
			response = "END Invalid choice."
		}
	default:
		response = "END Invalid choice or session timeout."
	}
	return c.SendString(response)
}

func (h *Handler) EventsHandler(c *fiber.Ctx) error {

	var durationInMillis = c.FormValue("durationInMillis", "")
	var phoneNumber = c.FormValue("phoneNumber", "")
	var errorMessage = c.FormValue("errorMessage", "")

	log.Println(fiber.Map{
		"durationInMillis": durationInMillis,
		"phoneNumber":      phoneNumber,
		"errorMessage":     errorMessage,
	})
	return nil
}

func (h *Handler) ImportKeystore(phoneNumber, pin string) (*accounts.Account, error) {
	ksFile, err := h.DB.SelectWalletByPhone(phoneNumber)
	if err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := os.ReadFile(ksFile.KeystorePath)
	if err != nil {
		fmt.Println("ks: ", err)
		return nil, err
	}

	account, err := ks.Import(jsonBytes, pin, pin)
	if err != nil {
		fmt.Println("acc: ", err)
		return nil, err
	}

	if err := os.Remove(ksFile.KeystorePath); err != nil {
		fmt.Println("osrem: ", err)
		return nil, err
	}
	return &account, nil
}
