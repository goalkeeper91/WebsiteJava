using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class GetSeason
    {
        public static string getActualSeason()
        {
            string season = null;
            string[] strings = null;
            string query = "Select cupname from cups ORDER BY id DESC LIMIT 1";
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;
            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                reader.Read();
                season = reader.GetString(0);
                strings = season.Split(' ');
                season = strings[0] + " " + strings[1] + " " + strings[2];
            }
            reader.Close();
            cnn.Close();

            return season;
        }
    }
}
