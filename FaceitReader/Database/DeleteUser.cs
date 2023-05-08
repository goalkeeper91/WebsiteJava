using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class DeleteUser
    {
        public static void deleteUser(int id)
        {
            MySqlConnection cnn = Connection.openConnection();
            MySqlCommand command;

            string query = "DELETE FROM user WHERE id = '" + id + "'";

            command = new MySqlCommand(query, cnn);
            command.ExecuteNonQuery();

            cnn.Close();
        }
    }
}
