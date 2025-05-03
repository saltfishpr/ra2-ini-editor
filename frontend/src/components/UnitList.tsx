import { Collapse, theme } from "antd";
import { main } from "../../wailsjs/go/models";
import "./UnitList.css";

interface UnitListProps {
  items: main.Unit[];
  onSelect: (unit: main.Unit) => void;
  selectedUnit: main.Unit | null; // 添加选中的单位属性
}

const UnitList = ({ items, onSelect, selectedUnit }: UnitListProps) => {
  const { token } = theme.useToken();

  // 按 type 对单位进行分组
  const groupedUnits = items.reduce(
    (acc: Record<string, main.Unit[]>, unit) => {
      const type = unit.type || "其他";
      if (!acc[type]) {
        acc[type] = [];
      }
      acc[type].push(unit);
      return acc;
    },
    {}
  );

  // 创建折叠面板的项目
  const collapseItems = Object.entries(groupedUnits).map(([type, units]) => ({
    key: type,
    label: type,
    children: (
      <div className="unit-list">
        {units.map((unit) => (
          <div
            key={unit.id}
            className={`unit-item ${
              selectedUnit && selectedUnit.id === unit.id ? "selected" : ""
            }`}
            onClick={() => onSelect(unit)}
          >
            {unit.ui_name || unit.name}
          </div>
        ))}
      </div>
    ),
  }));

  return (
    <Collapse
      accordion
      bordered={false}
      style={{ background: token.colorBgContainer }}
      items={collapseItems}
    />
  );
};

export default UnitList;
