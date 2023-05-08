using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class InsertInto
    {
        public static void CupPlaces(string place, string team, string cup)
        {
            string queryCheck = "SELECT cup_id FROM cupplaces WHERE cup_id = '" + cup + "'";
            MySqlConnection cnn2 = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            int counter = 0;
            command = new MySqlCommand(queryCheck, cnn2);
            reader = command.ExecuteReader();
            while (reader.Read())
            {
                counter++;
            }
            if (counter < 32)
            {
                string query = "INSERT INTO cupplaces (place,team_id,cup_id) VALUES ('" + place + "','" + team + "','" + cup + "')";
                MySqlConnection cnn = Connection.openConnection();
                command = new MySqlCommand(query, cnn);
                command.ExecuteNonQuery();
                cnn.Close();
            }
            reader.Close();
            cnn2.Close();
        }
    }
}
