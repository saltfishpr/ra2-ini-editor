import { PlusOutlined } from "@ant-design/icons";
import { AutoComplete, Button, Input, List, Tooltip, Typography } from "antd";
import { useEffect, useState } from "react";
import { ListAvailableProperties, NewULID } from "../../wailsjs/go/main/App";
import { main } from "../../wailsjs/go/models";

interface UnitDetailProps {
  value: main.Unit;
  onChange?: (value: main.Unit) => void;
}

const UnitDetail: React.FC<UnitDetailProps> = ({ value, onChange }) => {
  const [currentUnit, setCurrentUnit] = useState(value);
  const [availableProperties, setAvailableProperties] = useState<
    main.Property[]
  >([]);

  useEffect(() => {
    ListAvailableProperties(currentUnit.type)
      .then((props) => {
        setAvailableProperties(props);
      })
      .catch((error) => {
        console.error("Failed to fetch properties:", error);
      });
  }, [currentUnit.type]);

  const updateUnitProperties = (newProperties: main.Property[]) => {
    setCurrentUnit((prev) => {
      const updatedUnit = main.Unit.createFrom({
        ...prev,
        properties: newProperties,
      });
      if (onChange) {
        onChange(updatedUnit);
      }
      return updatedUnit;
    });
  };

  const handlePropertyChange = (
    index: number,
    field: keyof main.Property,
    value: string
  ) => {
    const newProperties = [...currentUnit.properties];
    newProperties[index] = { ...newProperties[index], [field]: value };
    updateUnitProperties(newProperties);
  };

  const handleAddProperty = () => {
    NewULID()
      .then((ukey) => {
        const newProperties = [
          ...currentUnit.properties,
          { ukey: ukey, key: "", value: "", comment: "" },
        ];
        updateUnitProperties(newProperties);
      })
      .catch((error) => {
        console.error("Failed to generate new ULID:", error);
      });
  };

  const handleDeleteProperty = (index: number) => {
    const newProperties = [...currentUnit.properties];
    newProperties.splice(index, 1);
    updateUnitProperties(newProperties);
  };

  return (
    <div style={{ padding: 24, border: "1px solid #d9d9d9", borderRadius: 8 }}>
      <Typography.Title level={4} style={{ marginBottom: 16 }}>
        基本信息
      </Typography.Title>
      <div
        style={{
          display: "flex",
          gap: 16,
          marginBottom: 24,
          padding: 16,
          borderRadius: 4,
        }}
      >
        <div style={{ flex: 1, display: "flex", alignItems: "center", gap: 8 }}>
          <Typography.Text style={{ whiteSpace: "nowrap" }}>
            Name:
          </Typography.Text>
          <Input value={currentUnit.name} disabled />
        </div>
        <div style={{ flex: 1, display: "flex", alignItems: "center", gap: 8 }}>
          <Typography.Text style={{ whiteSpace: "nowrap" }}>
            ID:
          </Typography.Text>
          <Input value={currentUnit.id} disabled />
        </div>
        <div style={{ flex: 1, display: "flex", alignItems: "center", gap: 8 }}>
          <Typography.Text style={{ whiteSpace: "nowrap" }}>
            Type:
          </Typography.Text>
          <Input value={currentUnit.type} disabled />
        </div>
      </div>

      <Typography.Title level={4} style={{ marginBottom: 16 }}>
        属性
      </Typography.Title>
      <List
        style={{ borderRadius: 4, padding: "0 16px" }}
        dataSource={currentUnit.properties}
        rowKey={(item) => item.ukey}
        renderItem={(item, index) => (
          <List.Item
            actions={[
              <Button
                type="text"
                danger
                onClick={() => handleDeleteProperty(index)}
              >
                删除
              </Button>,
            ]}
          >
            <div style={{ display: "flex", gap: 16, flex: 1 }}>
              <AutoComplete
                value={item.key}
                options={availableProperties.map((prop) => ({
                  value: prop.key,
                  label: (
                    <Tooltip
                      title={prop.desc}
                      placement="right"
                      style={{ width: "300px" }}
                    >
                      <div>{prop.key}</div>
                    </Tooltip>
                  ),
                }))}
                style={{ width: 200 }}
                filterOption={(inputValue, option) => {
                  if (!option) return false;
                  const optionValue = option.value.toLowerCase();
                  return optionValue.includes(inputValue.toLowerCase());
                }}
                onSelect={(value) => {
                  handlePropertyChange(index, "key", value);
                }}
                onChange={(value) => {
                  handlePropertyChange(index, "key", value);
                }}
              />
              <Input
                value={item.value}
                style={{ width: 200 }}
                placeholder="值"
                onChange={(e) =>
                  handlePropertyChange(index, "value", e.target.value)
                }
              />
              <Input
                value={item.comment}
                style={{ flex: 1 }}
                placeholder="注释"
                onChange={(e) =>
                  handlePropertyChange(index, "comment", e.target.value)
                }
              />
            </div>
          </List.Item>
        )}
      />
      <Button
        type="dashed"
        icon={<PlusOutlined />}
        style={{ width: "100%", margin: "12px 16px" }}
        onClick={handleAddProperty}
      >
        添加属性
      </Button>
    </div>
  );
};

export default UnitDetail;
