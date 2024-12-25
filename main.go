package main

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
	"sync"
)

// Funzione principale
func main() {
	config := readConfig()
	token, err := auth.RequestLiveToken()
	if err != nil {
		panic(err)
	}
	src := auth.RefreshTokenSource(token)

	// Configura il server proxy
	p, err := minecraft.NewForeignStatusProvider(config.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}
	listener, err := minecraft.ListenConfig{
		StatusProvider: p,
	}.Listen("raknet", config.Connection.LocalAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// Accetta connessioni in loop
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c.(*minecraft.Conn), listener, config, src)
	}
}

// Gestisce la connessione
func handleConn(conn *minecraft.Conn, listener *minecraft.Listener, config config, src oauth2.TokenSource) {
	serverConn, err := minecraft.Dialer{
		TokenSource: src,
		ClientData:  conn.ClientData(),
	}.Dial("raknet", config.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}

	var g sync.WaitGroup
	g.Add(2)

	// Invia messaggi iniziali al client
	go func() {
		if err := conn.StartGame(serverConn.GameData()); err != nil {
			panic(err)
		}

		// Messaggi di benvenuto
		conn.WritePacket(&packet.Text{Message: "§b[Ibrid]:§7 Ibrid Proxy running on " + protocol.CurrentVersion, TextType: packet.TextTypeChat})
		conn.WritePacket(&packet.Text{Message: "§b-----§a<AvailableCommands>§b-----", TextType: packet.TextTypeChat})
		conn.WritePacket(&packet.Text{Message: "§a/freeze §e> §bFreeze yourself bypassing walls", TextType: packet.TextTypeChat})
		conn.WritePacket(&packet.Text{Message: "§a/gamemode <int> §e> §bEnter in any gamemode", TextType: packet.TextTypeChat})
		conn.WritePacket(&packet.SetTitle{Text: "§b[IBRID]", ActionType: packet.TitleActionSetTitle})
		conn.WritePacket(&packet.SetTitle{Text: "§bdiscord.gg/4dCq3MbP", ActionType: packet.TitleActionSetSubtitle})
		conn.WritePacket(&packet.RequestPermissions{PermissionLevel: 1})

		// Log della connessione
		fmt.Println("[NetworkSession]->Client connected using Minecraft Version:", protocol.CurrentVersion, "[", protocol.CurrentProtocol, "]")
		fmt.Println("[NetworkSession]->Client found on:", config.Connection.RemoteAddress, "joined using ip:", config.Connection.LocalAddress, "or 127.0.0.1")
		g.Done()
	}()

	// Sincronizza con il server remoto
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		g.Done()
	}()
	g.Wait()

	// Riceve i pacchetti dal client
	go func() {
		defer listener.Disconnect(conn, "connection lost")
		defer fmt.Println("[NetworkSession]->Client Disconnected")
		defer serverConn.Close()

		for {
			pk, err := conn.ReadPacket()
			if err != nil {
				return
			}

			// Gestisce i comandi ricevuti
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/gamemode.0" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§a Gamemode§7 has been set to §aSurvival",
						TextType: packet.TextTypeChat,
					})
					conn.WritePacket(&packet.StartGame{PlayerGameMode: 0})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/gamemode.1" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§a Gamemode§7 has been set to §aCreative",
						TextType: packet.TextTypeChat,
					})
					conn.WritePacket(&packet.StartGame{PlayerGameMode: 1})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/gamemode.2" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§a Gamemode§7 has been set to §aSpectator",
						TextType: packet.TextTypeChat,
					})
					conn.WritePacket(&packet.StartGame{PlayerGameMode: 3})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/freeze" {
					// Send a message to the client
					fmt.Println("[NetworkSession]->Client executed | freeze | command")
					fmt.Println("[NetworkSession]->Client disconnected: PlayerNetworkSessionResponse = false")
					defer listener.Disconnect(conn, "§b[Ibrid]:§7 Generating too many packets (B)]")
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Generating packets (B)",
						TextType: packet.TextTypeChat,
					})

					for {

					}
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/crashA" {
					// Send a message to the client

					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/fly.true" {
					// Send a message to the clie
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/fly.false" {
					// Send a message to the clie
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/automine.diamond" {
					// Send a message to the clie
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/automine.gold" {
					// Send a message to the clie
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/automine.iron" {
					// Send a message to the clie
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if strings.HasPrefix(text.CommandLine, "/reachA") {
					for i, arg := range strings.Split(strings.TrimPrefix(text.CommandLine, "/reachA"), " ")[1:] {
						conn.WritePacket(&packet.Text{
							Message:  fmt.Sprintf("§b[Ibrid]:§7 Currently disabled(", arg, i),
							TextType: packet.TextTypeChat,
						})
					}
				}
				if strings.HasPrefix(text.CommandLine, "/reachA") {
					for i, arg := range strings.Split(strings.TrimPrefix(text.CommandLine, "/reachA"), " ")[2:] {
						conn.WritePacket(&packet.Text{
							Message:  fmt.Sprintf("§b[Ibrid]:§7 Invalid Arguments (%v)", i, arg),
							TextType: packet.TextTypeChat,
						})
					}
				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/autolog.true" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})
				}

			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/autolog.false" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
			}
			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/crashB" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
			}

			if text, ok := pk.(*packet.CommandRequest); ok {
				if text.CommandLine == "/autoclicker.1"+
					"/autoclicker.2"+"/autoclicker.3"+"/autoclicker.4"+"/autoclicker.5" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Cps should start from 6",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.7" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.8" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.9" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.10" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.11" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.12" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.13" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.14" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.15" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.16" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.17" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.18" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
				if text.CommandLine == "/autoclicker.19" {
					// Send a message to the client
					conn.WritePacket(&packet.Text{
						Message:  "§b[Ibrid]:§7 Currently disabled",
						TextType: packet.TextTypeChat,
					})

				}
			} else {
			}
			// Inoltra i pacchetti al server remoto
			if err := serverConn.WritePacket(pk); err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					_ = listener.Disconnect(conn, disc.Error())
				}
				return
			}
		}
	}()

	// Riceve i pacchetti dal server remoto
	go func() {
		defer serverConn.Close()
		defer listener.Disconnect(conn, "connection lost")

		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				var disc minecraft.DisconnectError
				if ok := errors.As(err, &disc); ok {
					_ = listener.Disconnect(conn, disc.Error())
				}
				return
			}
			if err := conn.WritePacket(pk); err != nil {
				return
			}
		}
	}()
}

// Configurazione
type config struct {
	Connection struct {
		LocalAddress  string
		RemoteAddress string
	}
}

// Legge la configurazione
func readConfig() config {
	c := config{}
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("create config: %v", err)
		}
		data, err := toml.Marshal(c)
		if err != nil {
			log.Fatalf("encode default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("write default config: %v", err)
		}
		_ = f.Close()
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		log.Fatalf("decode config: %v", err)
	}
	if c.Connection.LocalAddress == "" {
		c.Connection.LocalAddress = "0.0.0.0:19132"
	}
	data, _ = toml.Marshal(c)
	if err := os.WriteFile("config.toml", data, 0644); err != nil {
		log.Fatalf("write config: %v", err)
	}
	return c
}
