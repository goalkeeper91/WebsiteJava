using MySql.Data.MySqlClient;

namespace FaceitReader.Database
{
    internal class Connection
    {
        public static MySqlConnection openConnection()
        {
            MySqlConnection cnn;
            string conString = Globals.connectionString;
            cnn = new MySqlConnection(conString);
            cnn.Open();
            return cnn;
        }
    }
}
