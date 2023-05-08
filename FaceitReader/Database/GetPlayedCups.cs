using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using FaceitReader.Classes;
using MySql.Data.MySqlClient;

namespace FaceitReader.Database
{
    internal class GetPlayedCups
    {
        public static List<string> GetPlayedCupsLists(string teamId)
        {
            string query = "Select * from cups as c INNER JOIN places_cups as pc on pc.teams_id = '" + teamId + "' AND c.cupid = pc.cup_id";
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            List<string> playername = new List<string>();

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                while (reader.Read())
                {
                    playername.Add(reader.GetString(2));
                }
            }
            reader.Close();
            cnn.Close();

            return playername;
        }
    }
}
