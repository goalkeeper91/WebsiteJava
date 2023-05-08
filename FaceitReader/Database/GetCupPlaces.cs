using FaceitReader.Classes;
using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Data;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class GetCupPlaces
    {
        public static List<FullPlaceList> getSelectedPlaceList(string cupId)
        {
            List<FullPlaceList> result = new List<FullPlaceList>();
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            string query;
            string price = null;

            if (!string.IsNullOrWhiteSpace(cupId))
            {
                query = "SELECT * FROM cupplaces AS cp LEFT JOIN teams AS t ON cp.cup_id = '" + cupId + "' AND t.team_id = cp.team_id LEFT JOIN player AS p ON t.captain_id = p.player_id ORDER BY cp.place";
                command = new MySqlCommand(query, cnn);
                reader = command.ExecuteReader();
                while (reader.Read())
                {
                    if(Convert.ToInt32(reader.GetString(0)) == 1)
                    {
                        price = "25 Punkte + 5x 25€";
                    }
                    else if (Convert.ToInt32(reader.GetString(0)) == 2)
                    {
                        price = "20 Punkte + 5x 15€";
                    }
                    else if (Convert.ToInt32(reader.GetString(0)) > 2 && Convert.ToInt32(reader.GetString(0)) < 5)
                    {
                        price = "15 Punkte + 5x 10€";
                    }
                    else if (Convert.ToInt32(reader.GetString(0)) > 4 && Convert.ToInt32(reader.GetString(0)) < 9)
                    {
                        price = "10 Punkte + 5x 5€";
                    }
                    else if (Convert.ToInt32(reader.GetString(0)) == 6)
                    {
                        price = "5x 2€";
                    }
                    result.Add(new FullPlaceList() { Place = Convert.ToInt32(reader.GetString(1)), TeamName = reader.GetString(7), TeamId = reader.GetString(2), CaptainName = reader.GetString(15), Prize = price });
                }
                reader.Close();
            }
            cnn.Close();
            Console.WriteLine(result.Count);
            return result;
        }
        public static List<FullPlaceList> getCup(string season, string cup)
        {
            string cupId = "";
            string query = "SELECT cupid FROM cups WHERE cupname = 'SKINBARON Cup " + season + " Cup " + cup + "'";
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            List<FullPlaceList> data;

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();

            if (reader.HasRows)
            {
                reader.Read();
                cupId = reader.GetValue(0).ToString();
            }
            
            reader.Close();
            cnn.Close();

            data = getSelectedPlaceList(cupId);

            return data;
        }
        public static List<string> getCupName()
        {
            string query = "SELECT cupname FROM cups";
            List<string> cupname = new List<string>();
            string name;
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                while (reader.Read())
                {
                    name = reader.GetValue(0).ToString();
                    name = name.Replace("SKINBARON","").Replace("Cup","").Replace(" ", ",");
                    cupname.Add(name);
                }
            }

            reader.Close();
            cnn.Close();

            return cupname;
        }
    }
}
