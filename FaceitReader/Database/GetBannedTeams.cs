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
            if (!reader.HasRows)
            {
                name = reader.GetValue(1).ToString();
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
