using FaceitReader.Classes;
using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;

namespace FaceitReader.Database
{
    internal class GetUser
    {
        public static List<User> getUser()
        {
            List<User> users = new List<User>();
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            string firstname = null;
            string lastname = null;
            string isAdmin;

            string query = "Select * from user";

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                while (reader.Read())
                {
                    if(!reader.IsDBNull(4))
                    {
                        firstname = reader.GetString(4);
                    }
                    if(!reader.IsDBNull(5))
                    {
                        lastname = reader.GetString(5);
                    }
                    Console.WriteLine(reader.GetInt32(6));
                    if(reader.GetInt32(6) == 1)
                    {
                        isAdmin = "Ja";
                    }
                    else
                    {
                        isAdmin = "Nein";
                    }
                    users.Add(new User() { Id = reader.GetInt32(0), Username = reader.GetString(1), Email = reader.GetString(2), Firstname = firstname, Lastname = lastname, Admin = isAdmin });
                }
            }
            reader.Close();
            cnn.Close();

            return users;
        }
    }
}
