using System;
using System.Runtime.Remoting.Channels;
using System.Windows.Forms;

namespace FaceitReader
{
    static class Program
    {
        /// <summary>
        /// Der Haupteinstiegspunkt für die Anwendung.
        /// </summary>
        [STAThread]
        static void Main()
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new PlaceOverview());
        }
    }
}
