import { type LucideIcon, StarIcon } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useCallback, useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import ReviewSection from "@/components/Libro/ReviewSection";
import SectionMenuItem from "@/components/SectionMenuItem";
import { useTranslate } from "@/utils/i18n";

type SettingSection = "review";

interface State {
  selectedSection: SettingSection;
}

const BASIC_SECTIONS: SettingSection[] = ["review"];
const SECTION_ICON_MAP: Record<SettingSection, LucideIcon> = {
  review: StarIcon,
};

const Libro = observer(() => {
  const t = useTranslate();
  const location = useLocation();
  const [state, setState] = useState<State>({
    selectedSection: "review",
  });

  const settingsSectionList = [...BASIC_SECTIONS];

  useEffect(() => {
    let hash = location.hash.slice(1) as SettingSection;
    // If the hash is not a valid section, redirect to the default section.
    if (![...BASIC_SECTIONS].includes(hash)) {
      hash = "review";
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
                text={t(`libro.${item}`)}
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
                    {t(`libro.${settingSection}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {state.selectedSection === "review" ? (
            <ReviewSection />
          ) : null}
        </div>
      </div>
    </section>
  )
});

export default Libro;
