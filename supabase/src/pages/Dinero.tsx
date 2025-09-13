import { CalendarDaysIcon, ChartBarIcon, ListIcon, type LucideIcon, WalletCardsIcon } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useCallback, useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import CategorySection from "@/components/Dinero/CategorySection";
import CalendarSection from "@/components/Dinero/CalendarSection";
import ChartByCategorySection from "@/components/Dinero/ChartByCategorySection";
import ExpenseSection from "@/components/Dinero/ExpenseSection";
import SectionMenuItem from "@/components/SectionMenuItem";
import { useTranslate } from "@/utils/i18n";

type SettingSection = "category" | "expense" | "by-category" | "by-calendar";

interface State {
  selectedSection: SettingSection;
}

const BASIC_SECTIONS: SettingSection[] = ["expense", "category"];
const CHART_SECTIONS: SettingSection[] = ["by-calendar", "by-category"];
const SECTION_ICON_MAP: Record<SettingSection, LucideIcon> = {
  expense: WalletCardsIcon,
  category: ListIcon,
  "by-calendar": CalendarDaysIcon,
  "by-category": ChartBarIcon,
};

const Dinero = observer(() => {
  const t = useTranslate();
  const location = useLocation();
  const [state, setState] = useState<State>({
    selectedSection: "expense",
  });

  const settingsSectionList = [...BASIC_SECTIONS, ...CHART_SECTIONS];

  useEffect(() => {
    let hash = location.hash.slice(1) as SettingSection;
    // If the hash is not a valid section, redirect to the default section.
    if (![...BASIC_SECTIONS, ...CHART_SECTIONS].includes(hash)) {
      hash = "expense";
    }

    setState({
      selectedSection: hash,
    });
  }, [location.hash]);

  const handleSectionSelectorItemClick = useCallback((settingSection: SettingSection) => {
    window.location.hash = settingSection;
  }, []);

  return (
    <section className="@container w-full max-w-6xl min-h-full flex flex-col justify-start items-center pt-6 pb-8">
      <div className="w-full border border-border flex flex-col flex-row justify-start items-start px-4 py-3 rounded-xl bg-background text-foreground">
        <div className="flex flex-col justify-start items-start w-40 h-auto shrink-0 py-2">
          <span className="text-sm mt-0.5 pl-3 font-mono select-none text-muted-foreground">{t("common.basic")}</span>
          <div className="w-full flex flex-col justify-start items-start mt-1">
            {BASIC_SECTIONS.map((item) => (
              <SectionMenuItem
                key={item}
                text={t(`dinero.${item}`)}
                icon={SECTION_ICON_MAP[item]}
                isSelected={state.selectedSection === item}
                onClick={() => handleSectionSelectorItemClick(item)}
              />
            ))}
          </div>
          <span className="text-sm mt-4 pl-3 font-mono select-none text-muted-foreground">{t("common.chart")}</span>
          <div className="w-full flex flex-col justify-start items-start mt-1">
            {CHART_SECTIONS.map((item) => (
              <SectionMenuItem
                key={item}
                text={t(`dinero.${item}`)}
                icon={SECTION_ICON_MAP[item]}
                isSelected={state.selectedSection === item}
                onClick={() => handleSectionSelectorItemClick(item)}
              />
            ))}
          </div>
        </div>
        <div className="w-full grow pl-4 overflow-x-auto">
          <div className="w-auto my-2 hidden">
            <Select value={state.selectedSection} onValueChange={(value) => handleSectionSelectorItemClick(value as SettingSection)}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select section" />
              </SelectTrigger>
              <SelectContent>
                {settingsSectionList.map((settingSection) => (
                  <SelectItem key={settingSection} value={settingSection}>
                    {t(`dinero.${settingSection}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {state.selectedSection === "expense" ? (
            <ExpenseSection />
          ) : state.selectedSection === "category" ? (
            <CategorySection />
          ) : state.selectedSection === "by-calendar" ? (
            <CalendarSection />
          ) : state.selectedSection === "by-category" ? (
            <ChartByCategorySection />
          ) : null}
        </div>
      </div>
    </section>
  )
});

export default Dinero;
