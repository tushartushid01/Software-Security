export const ButtonList = [
  "Electronics",
  "Mobile",
  "TV",
  "Appliances",
  "Fashion",
  "Women-Fashion",
];

export enum Button {
  mobile = "Mobile",
  tv = "TV",
  appliances = "Appliances",
  fashion = "Fashion",
  "women-fashion" = "Women-Fashion",
}

export interface IData {
  id: string;
  title: string;
  desription: string;
}

export const dataDummy = [
  {
    id: "1",
    title: "Redmi 12 ",
    desription: "5G Moonstone Silver 6GB RAM 128GB ROM",
  },
  {
    id: "2",
    title: "POCO C51 ",
    desription: "(Royal Blue, 6GB RAM, 128GB Storage)",
  },
  {
    id: "3",
    title: "POCO C51 ",
    desription: "(Power Black, 6GB RAM, 128GB Storage)",
  },
  {
    id: "4",
    title: "Redmi 13C",
    desription:
      " (Starfrost White, 4GB RAM, 128GB Storage) | Powered by 4G MediaTek Helio G85 | 90Hz Display | 50MP AI Triple Camera",
  },
  {
    id: "1",
    title: "beatXP Kitchen ",
    desription:
      "beatXP Kitchen Scale Multipurpose Portable Electronic Digital Weighing",
  },
  {
    id: "2",
    title: "GLUN®",
    desription:
      " Electronic Portable Digital LED Screen Luggage Weighing Scale, 50 kg/110 Lb",
  },
  {
    id: "1",
    title: "Russell Hobbs ",
    description:
      "Russell Hobbs England RFS800-800 Watt Food Steamer Steam Cooker with 2 Years Manufacturer Warranty (White)",
  },
  {
    id: "2",
    title: "Microwave",
    description:
      "SHARP 27L Convection Microwave Oven (R827KNK, Black, Healthy Fry with Zero Oil, Stainless Steel Cavity, Vapour Clean)",
  },
  {
    id: "3",
    title: "Air Fryer",
    description:
      "Pigeon Healthifry Digital Air Fryer, 360° High Speed Air Circulation Technology 1200 W with Non-Stick 4.2 L Basket - Green",
  },
];
