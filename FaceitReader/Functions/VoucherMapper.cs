using FaceitReader.Classes;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Functions
{
    public class VoucherMapper
    {
        public static List<Voucher> MapFromRangeData(IList<IList<object>> values, int id)
        {
            var items = new List<Voucher>();
            foreach (var value in values)
            {
                if(id == 1)
                {
                    Voucher item = new Voucher()
                    {
                        Price1 = value[0].ToString(),
                        Price2 = value[1].ToString(),
                        Price3 = value[2].ToString(),
                        Price4 = value[3].ToString(),
                        Price5 = value[4].ToString()
                    };
                    items.Add(item);
                }
                else
                {
                    Voucher item = new Voucher()
                    {
                        Price1 = value[0].ToString(),
                        Price3 = value[1].ToString(),
                        Price4 = value[2].ToString(),
                        Price5 = value[3].ToString(),
                    };
                    items.Add(item);
                }
            }
            return items;
        }
        public static IList<IList<object>> MapToRangeData(Voucher item)
        {
            var objectList = new List<object>() { item.Price1, item.Price2, item.Price3, item.Price4, item.Price5 };
            var rangeData = new List<IList<object>> { objectList };
            return rangeData;
        }
    }
}
