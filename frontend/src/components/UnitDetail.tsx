import { PlusOutlined } from "@ant-design/icons";
import { Button, Input, List, Typography } from "antd";
import { useState } from "react";
import { main } from "../../wailsjs/go/models";

interface UnitDetailProps {
  unit: main.Unit;
  onSave: (unit: main.Unit) => void;
}

const UnitDetail = ({ unit, onSave }: UnitDetailProps) => {
  const [localUnit, setLocalUnit] = useState(unit);

  return (
    <div
      style={{
        padding: 24,
        border: "1px solid #d9d9d9",
        borderRadius: 8,
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "flex-end",
          gap: 8,
          marginBottom: 16,
        }}
      >
        <Button>Copy</Button>
        <Button
          type="primary"
          onClick={() => {
            onSave(localUnit);
          }}
        >
          Save
        </Button>
      </div>
      <div
        style={{
          display: "flex",
          gap: 16,
          marginBottom: 24,
          padding: 16,
          borderRadius: 4,
        }}
      >
        <div style={{ flex: 1 }}>
          <Typography.Text>Name:</Typography.Text>
          <Input
            value={localUnit.name}
            onChange={(e) =>
              setLocalUnit(
                main.Unit.createFrom({
                  ...localUnit,
                  name: e.target.value,
                })
              )
            }
          />
        </div>
        <div style={{ flex: 1 }}>
          <Typography.Text>ID:</Typography.Text>
          <Input
            value={localUnit.id}
            onChange={(e) =>
              setLocalUnit(
                main.Unit.createFrom({
                  ...localUnit,
                  id: parseInt(e.target.value),
                })
              )
            }
          />
        </div>
        <div style={{ flex: 1 }}>
          <Typography.Text>Type:</Typography.Text>
          <Input
            value={localUnit.type}
            onChange={(e) =>
              setLocalUnit(
                main.Unit.createFrom({
                  ...localUnit,
                  type: e.target.value,
                })
              )
            }
          />
        </div>
      </div>

      <Typography.Title level={4} style={{ marginBottom: 16 }}>
        属性
      </Typography.Title>

      <List
        style={{
          borderRadius: 4,
          padding: "0 16px",
        }}
        dataSource={localUnit.properties}
        renderItem={(item, index) => (
          <List.Item
            actions={[
              <Button type="text" danger>
                删除
              </Button>,
            ]}
          >
            <div style={{ display: "flex", gap: 16, flex: 1 }}>
              <Input
                value={item.key}
                style={{ width: 200 }}
                placeholder="键"
                onChange={(e) => {
                  const newProperties = [...localUnit.properties];
                  newProperties[index].key = e.target.value;
                  setLocalUnit(
                    main.Unit.createFrom({
                      ...localUnit,
                      properties: newProperties,
                    })
                  );
                }}
              />
              <Input
                value={item.value}
                style={{ width: 200 }}
                placeholder="值"
                onChange={(e) => {
                  const newProperties = [...localUnit.properties];
                  newProperties[index].value = e.target.value;
                  setLocalUnit(
                    main.Unit.createFrom({
                      ...localUnit,
                      properties: newProperties,
                    })
                  );
                }}
              />
              <Input
                value={item.comment}
                style={{ flex: 1 }}
                placeholder="注释"
                onChange={(e) => {
                  const newProperties = [...localUnit.properties];
                  newProperties[index].comment = e.target.value;
                  setLocalUnit(
                    main.Unit.createFrom({
                      ...localUnit,
                      properties: newProperties,
                    })
                  );
                }}
              />
            </div>
          </List.Item>
        )}
      />

      <Button
        type="dashed"
        icon={<PlusOutlined />}
        style={{
          width: "100%",
          margin: "12px 16px",
        }}
        onClick={() => {
          setLocalUnit(
            main.Unit.createFrom({
              ...localUnit,
              properties: [
                ...localUnit.properties,
                { key: "", value: "", comment: "" },
              ],
            })
          );
        }}
      >
        添加属性
      </Button>
    </div>
  );
};

export default UnitDetail;
