using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class GetPlayer
    {
        public static List<string> getPlayer(string teamid)
        {
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            List<string> playername = new List<string>();

            string query = "Select * from player where team_id = '" + teamid + "'";

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
                if (reader.HasRows)
                {
                    while (reader.Read())
                    {
                        playername.Add(reader.GetString(5));
                    }
                }
                reader.Close();
                cnn.Close();

            return playername;
        }
    }
}
