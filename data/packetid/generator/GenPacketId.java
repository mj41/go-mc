package pers.tnze.gomc.gen;

import net.minecraft.SharedConstants;
import net.minecraft.network.ProtocolInfo;
import net.minecraft.network.protocol.configuration.ConfigurationProtocols;
import net.minecraft.network.protocol.game.GameProtocols;
import net.minecraft.network.protocol.login.LoginProtocols;
import net.minecraft.network.protocol.status.StatusProtocols;
import net.minecraft.server.Bootstrap;

import java.io.FileWriter;
import java.io.IOException;
import java.io.Writer;

// 1.21.11: ProtocolInfo.Unbound → UnboundProtocol/SimpleUnboundProtocol (both extend DetailsProvider)

public class GenPacketId {

    public static void main(String[] args) throws Exception {
        SharedConstants.tryDetectVersion();
        Bootstrap.bootStrap();
        try (FileWriter w = new FileWriter("packet_names.txt")) {
            handlePackets(w, LoginProtocols.CLIENTBOUND_TEMPLATE.details(), "ClientboundLogin");
            handlePackets(w, LoginProtocols.SERVERBOUND_TEMPLATE.details(), "ServerboundLogin");
            handlePackets(w, StatusProtocols.CLIENTBOUND_TEMPLATE.details(), "ClientboundStatus");
            handlePackets(w, StatusProtocols.SERVERBOUND_TEMPLATE.details(), "ServerboundStatus");
            handlePackets(w, ConfigurationProtocols.CLIENTBOUND_TEMPLATE.details(), "ClientboundConfig");
            handlePackets(w, ConfigurationProtocols.SERVERBOUND_TEMPLATE.details(), "ServerboundConfig");
            handlePackets(w, GameProtocols.CLIENTBOUND_TEMPLATE.details(), "");
            handlePackets(w, GameProtocols.SERVERBOUND_TEMPLATE.details(), "");
        }
    }

    private static void handlePackets(Writer w, ProtocolInfo.Details packets, String prefix) throws IOException {
        packets.listPackets((packetType, i) -> {
            String packetName = packetType.id().getPath();
            String[] words = packetName.split("_");
            try {
                if (prefix != null) {
                    w.write(prefix);
                }
                for (String word : words) {
                    w.write(Character.toUpperCase(word.charAt(0)));
                    w.write(word.substring(1));
                }
                w.write("\n");
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });
        w.write('\n');
    }
}