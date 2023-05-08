using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Classes
{
    public class Voucher
    {
        public int Id { get; set; }
        public string Price1 { get; set; }
        public string Price2 { get; set; }
        public string Price3 { get; set; }
        public string Price4 { get; set; }
        public string Price5 { get; set; }

        public override string ToString()
        {
            if (Price5 == null)
            {
                return "25€: " + Price1 + "  10€: " + Price2 + " 5€: " + Price3 + "  2€: " + Price4 ;
            }else
            {
                return "25€: " + Price1 + "  15€: " + Price2 + " 10€: " + Price3 + "  5€: " + Price4 + " 2€: " + Price5;
            }
        }
        public override bool Equals(object obj)
        {
            if (obj == null) return false;
            Voucher objAsFullPlaceList = obj as Voucher;
            if (objAsFullPlaceList == null) return false;
            else return Equals(objAsFullPlaceList);
        }
        public override int GetHashCode()
        {
            return Id;
        }
        public bool Equals(Voucher other)
        {
            if (other == null) return false;
            return (this.Id.Equals(other.Id));
        }
    }
}

