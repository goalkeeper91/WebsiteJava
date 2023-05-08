using MySql.Data.MySqlClient;

namespace FaceitReader.Database
{
    internal class GetBannedTeams
    {
        public static string getBannedTeam(string id)
        {
            string name;
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;

            string query = "Select * from banned_teams where team_id ='" + id+ "'";

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                reader.Read();
                name = reader.GetValue(2).ToString();
            }
            else
            {
                name = "";
            }
            reader.Close();
            cnn.Close();

            return name;
        }

        public static string getBannedPlayer(string id, string teamId)
        {
            string name;
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader;

            string query = "Select * from banned_players where player_id ='" + id + "' and not team_id = '" + teamId + "'";

            command = new MySqlCommand(query, cnn);
            reader = command.ExecuteReader();
            if (reader.HasRows)
            {
                reader.Read();
                name = reader.GetValue(2).ToString();
            }
            else
            {
                name = "";
            }
            reader.Close();
            cnn.Close();

            return name;
        }
    }
}
