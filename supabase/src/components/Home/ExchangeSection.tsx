import { observer } from "mobx-react-lite";
import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { type Exchange, homeStore } from "@/types/model/home_service";
import { toCamelCase } from "@/utils/common";

const ExchangeSection = observer(() => {
  const [exchanges, setExchanges] = useState<Exchange[]>([]);

  useEffect(() => {
    fetchExchanges();
  }, []);

  const fetchExchanges = async () => {
    try {
      let { data, error } = await homeStore.fetchExchanges();
      if (error != null) {
        throw error;
      }

      if (data) {
        setExchanges(data.map((e) => (toCamelCase(e))))
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.message);
    }
  };

  return (
    <div className="border border-dashed rounded-md flex flex-col justify-start items-start mt-4 px-2 py-2">
      {exchanges.map((exchange) => (
        <div key={exchange.id} className="w-full text-sm flex flex-col flex-row gap-2 px-2 justify-between items-center">
          <p className="font-small text-muted-foreground">1 {exchange.from}</p>
          <div >{exchange.value} {exchange.to}</div>
        </div>
      ))}
    </div>
  );
});

export default ExchangeSection;
