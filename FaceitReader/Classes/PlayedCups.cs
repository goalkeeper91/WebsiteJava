using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Classes
{
    internal class PlayedCupsList
    {
        public int Id { get; set; }
        public string Cup { get; set; }
        public override string ToString()
        {
            return "Cup: " + Cup;
        }
        public override bool Equals(object obj)
        {
            if (obj == null) return false;
            PlayedCupsList objAsPlayedCupsList = obj as PlayedCupsList;
            if (objAsPlayedCupsList == null) return false;
            else return Equals(objAsPlayedCupsList);
        }
        public override int GetHashCode()
        {
            return Id;
        }
        public bool Equals(PlayedCupsList other)
        {
            if (other == null) return false;
            return (this.Id.Equals(other.Id));
        }
    }
}
