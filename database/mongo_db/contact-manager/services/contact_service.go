package main

import (
	"fmt"
	"log"
	"os"

	"code_test/message_queue/mongodb_database/contact-manager/config"
	"code_test/message_queue/mongodb_database/contact-manager/services"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var contactService *services.ContactService

func init() {
	if err := config.InitDb(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	contactService = services.NewContactService()
}

func main() {
	var rootCmd = &cobra.Command{Use: "contact-manager"}

	rootCmd.AddCommand(createCmd())
	rootCmd.AddCommand(getByIDCmd())
	rootCmd.AddCommand(getByNameCmd())
	rootCmd.AddCommand(updateCmd())
	rootCmd.AddCommand(deleteCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config.CloseDB()
}

func createCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [name] [phone] [email]",
		Short: "Create a new contact",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			name, phone, email := args[0], args[1], args[2]
			contact, err := contactService.CreateContact(name, phone, email)
			if err != nil {
				fmt.Printf("Failed to create contact: %v\n", err)
				return
			}
			fmt.Printf("Created contact: %+v\n", contact)
		},
	}
}

func getByIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-id [id]",
		Short: "Get contact by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := primitive.ObjectIDFromHex(args[0])
			if err != nil {
				fmt.Printf("Invalid ID format: %v\n", err)
				return
			}
			contact, err := contactService.GetContactByID(id)
			if err != nil {
				fmt.Printf("Failed to get contact: %v\n", err)
				return
			}
			fmt.Printf("Found contact: %+v\n", contact)
		},
	}
}

func getByNameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-name [name]",
		Short: "Get contacts by name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			contacts, err := contactService.GetContactsByName(args[0])
			if err != nil {
				fmt.Printf("Failed to get contacts: %v\n", err)
				return
			}
			for _, contact := range contacts {
				fmt.Printf("Found contact: %+v\n", contact)
			}
		},
	}
}

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [id] [name] [phone] [email]",
		Short: "Update a contact",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := primitive.ObjectIDFromHex(args[0])
			if err != nil {
				fmt.Printf("Invalid ID format: %v\n", err)
				return
			}
			name, phone, email := args[1], args[2], args[3]
			if err := contactService.UpdateContact(id, name, phone, email); err != nil {
				fmt.Printf("Failed to update contact: %v\n", err)
				return
			}
			fmt.Println("Contact updated successfully")
		},
	}
}

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a contact",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := primitive.ObjectIDFromHex(args[0])
			if err != nil {
				fmt.Printf("Invalid ID format: %v\n", err)
				return
			}
			if err := contactService.DeleteContact(id); err != nil {
				fmt.Printf("Failed to delete contact: %v\n", err)
				return
			}
			fmt.Println("Contact deleted successfully")
		},
	}
}
